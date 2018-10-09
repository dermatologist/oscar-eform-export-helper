package main

import (
	"database/sql"
	"flag"
	"fmt"
	"github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/agent"
	"net"
	"os"
	"strconv"
)

type ViaSSHDialer struct {
	client *ssh.Client
}

func (self *ViaSSHDialer) Dial(addr string) (net.Conn, error) {
	return self.client.Dial("tcp", addr)
}

func mysqlConnect() (*sql.Rows, error) {
	sshHost := flag.String("host", "example.com", "The SSH host")
	sshPort := flag.Int("port", 22, "The port number")
	sshUser := flag.String("sshuser", "ssh-user", "ssh user")
	sshPass := flag.String("sshpass", "ssh-pass", "SSH Password")
	dbUser := flag.String("dbuser", "dbuser", "The db user")
	dbPass := flag.String("dbpass", "dbpass", "The db password")
	dbHost := flag.String("dbhost", "localhost:3306", "The db host")
	dbName := flag.String("dbname", "oscar", "The database name")
	fid := flag.Int("fid", 1, "The eform ID")

	var agentClient agent.Agent
	// Establish a connection to the local ssh-agent
	if conn, err := net.Dial("unix", os.Getenv("SSH_AUTH_SOCK")); err == nil {
		defer conn.Close()

		// Create a new instance of the ssh agent
		agentClient = agent.NewClient(conn)
	}

	// The client configuration with configuration option to use the ssh-agent
	sshConfig := &ssh.ClientConfig{
		User: *sshUser,
		Auth: []ssh.AuthMethod{},
	}

	// When the agentClient connection succeeded, add them as AuthMethod
	if agentClient != nil {
		sshConfig.Auth = append(sshConfig.Auth, ssh.PublicKeysCallback(agentClient.Signers))
	}
	// When there's a non empty password add the password AuthMethod
	if *sshPass != "" {
		sshConfig.Auth = append(sshConfig.Auth, ssh.PasswordCallback(func() (string, error) {
			return *sshPass, nil
		}))
	}

	// Connect to the SSH Server
	if sshcon, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", *sshHost, *sshPort), sshConfig); err == nil {
		defer sshcon.Close()

		// Now we register the ViaSSHDialer with the ssh connection as a parameter
		mysql.RegisterDial("mysql+tcp", (&ViaSSHDialer{sshcon}).Dial)

		// And now we can use our new driver with the regular mysql connection string tunneled through the SSH connection
		if db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@mysql+tcp(%s)/%s", *dbUser, *dbPass, *dbHost, *dbName)); err == nil {

			fmt.Printf("Successfully connected to the db\n")

			sqlQuery := `
			SELECT id, fdid, fid, demographic_no, var_name, 
			var_value FROM eform_values WHERE fid = " + strconv.Itoa(*fid) + " AND fdid IN
			 (SELECT fdid FROM eform_data WHERE fid = " + strconv.Itoa(*fid) + " AND form_date >= " 
			 + *dateFrom + " AND form_date <= " + *dateTo + " );
			 `

			if rows, err := db.Query(sqlQuery); err == nil {
				//for rows.Next() {
				//	var id int64
				//	var name string
				//	rows.Scan(&id, &name)
				//	fmt.Printf("ID: %d  Name: %s\n", id, name)
				//}
				return rows, nil
				defer rows.Close()
			} else {
				fmt.Printf("Failure: %s", err.Error())
				return nil, err
			}

			defer db.Close()

		} else {

			fmt.Printf("Failed to connect to the db: %s\n", err.Error())
			return nil, err
		}

	}
	return nil, nil
}
