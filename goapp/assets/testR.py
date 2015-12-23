#!/usr/bin/env python
from __future__ import print_function
import sys
import requests
import os
import re
import string
import subprocess

host = "http://testr-1154.appspot.com"

def header():
	return {'Content-Type' : 'application/json', 'Authorization':os.environ.get('TESTR_TOKEN')} 

def repo():
	return os.environ.get('TRAVIS_REPO_SLUG')	

def api():
	return host + "/api/repos/" + repo() +"/"

def home():
	return host + "/repos/" + repo() +"/"


def pushReport(c, r):
	req = {
		"Commit" : c,
		"Content": r,
		"Type":"surefire"
	}	
	res = requests.post(api(), json=req, headers=header())	
	if res.status_code == 201:
		return True
	else:
		print("ERROR %d\n:%s" % (res.status_code, res.text), file=sys.stderr)
		return False
	

def scanReports(r):	
	print("Scanning %s" % r)
	reports = []
	dirList = os.listdir(r)
	p = r + '/target/surefire-reports/TEST-TestSuite.xml'

	if os.path.os.access(p, os.R_OK):		
		reports = reports + [readReport(p)]
	for dir in dirList:				
		if os.path.isdir(dir) == True:						
			reports = reports + scanReports(dir)			
	return reports

def readReport(p):		
	f = open(p, 'r')
	cnt=""
	while True:	
		line = f.readline()
		if not line: break
		cnt += line		
	return cnt

def usage():
	print("testr.py [--host endpoint]")
	print("Default endpoint: " + host)
	print("Environment variables:")
	print("\t$TESTR_TOKEN: the testr token used to authenticate your requests")
	print("\t$GH_TOKEN: the github token to use to query the GitHub API")
	
	print("\nThe following variables shall be provided by Travis:")
	print("\t$TRAVIS_COMMIT: the SHA1 commit associated to the test reports")
	print("\t$TRAVIS_REPO_SLUG: the github repo. Formatted as 'owner/repository'")	

def getGhCommit(sha):	
	hdr = {'Content-Type' : 'application/json', 'Authorization': 'token %s' % os.environ.get('GH_TOKEN')} 


	out=requests.get("https://api.github.com/repos/%s/commits/%s" % (repo(), sha), headers=hdr)
	if out.status_code != 200:
		print("ERROR %s:\n%s" % (out.status_code, out.text))
	msg = out.json()	
	return {
		"Sha1": msg["sha"],
		"Url": msg["url"],
		"Log": msg["commit"]["message"],
		"Committer": {
			"Name": msg["commit"]["committer"]["name"],
			"Email": msg["commit"]["committer"]["email"],
			"Date": msg["commit"]["committer"]["date"],
			"AvatarUrl": msg["committer"]["avatar_url"]
		},
		"Author": {
			"Name": msg["commit"]["author"]["name"],
			"Email": msg["commit"]["author"]["email"],
			"Date": msg["commit"]["author"]["date"],
			"AvatarUrl": msg["author"]["avatar_url"]
		},		
	}

####### ---------- MAIN ------------- ################
if __name__ == "__main__":

	if len(sys.argv) == 2 and (sys.argv[1] == "-h" or sys.argv[1] == "--help"):
		usage()
		exit(0)

	sha = os.environ.get('TRAVIS_COMMIT')

	if not os.environ.get('GH_TOKEN'):
		print("Environment variable GH_TOKEN missing", file=sys.stderr)
		exit(1)

	if not os.environ.get('TRAVIS_REPO_SLUG'):
		print("Environment variable TRAVIS_REPO_SLUG missing", file=sys.stderr)
		exit(1)

	if not os.environ.get('TESTR_TOKEN'):
		print("Environment variable TESTR_TOKEN missing", file=sys.stderr)
		exit(1)

	if not sha and len(sys.argv) == 1:
		usage()				

	if len(sys.argv) == 3 and sys.argv[1] == "--host":
		host = sys.argv[2]

	root = "."
	if not os.path.os.access(root + "/pom.xml", os.R_OK):
		print("The script must be executed inside a maven project")
		exit(1)
	reports = scanReports(root)	
	if len(reports) == 0:
		print("No surefire reports available", file=sys.stderr)
		exit(0)				
	c = getGhCommit(sha)		
	done = False
	for r in reports:
		if pushReport(c, r):
			print("+", end="")
			done = done or True
		else:
			print("-", end="")
			done = done or False
	print()	
	if done:
		print("Reports uploaded. Open %s to see the reports" % home())