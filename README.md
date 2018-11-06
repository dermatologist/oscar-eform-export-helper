# :bookmark_tabs: GoForm - eForm Reporting Tool 
[![Build Status](https://travis-ci.org/dermatologist/oscar-eform-export-helper.svg?branch=develop)](https://travis-ci.org/dermatologist/oscar-eform-export-helper)

[![eForm Export Tool](https://raw.github.com/dermatologist/oscar-eform-export-helper/develop/notes/usage.gif)](https://canehealth.com)

## About

Reporting from OSCAR eforms is difficult as the fields are added to a single table as key value pairs. This is a tool for quickly diplaying summary statistics by connecting to the database via SSH. It can also work on the csv file exported from the eForm Export Tool.

## How to use

```
oscar_helper -sshhost=xxx -sshport=22 -sshuser=xxx -sshpass=xxx -dbuser=xxx -dbpass=xxx -dbname=xxx -dbhost=localhost -datefrom=YYYY-MM-DD -dateto=YYYY-MM-DD -fid=1 -include

```

## Work in progress: DO NOT USE

## Check out the :fire: [FHIRForm framework for E-Form management](https://github.com/E-Health/fhirform)

## Contributors

[Bell Eapen](http://canehealth.com)
