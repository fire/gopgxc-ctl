/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. 
 * Copyright (c) 2014, K. S. Ernest "iFire" Lee */
 
package exec

import (
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/agent"
	"io/ioutil"
	"net"
	"os"
)

func authAgent(conn net.Conn) (auth ssh.AuthMethod) {
	ag := agent.NewClient(conn)
	return ssh.PublicKeysCallback(ag.Signers)
}

// http://kukuruku.co/hub/golang/ssh-commands-execution-on-hundreds-of-servers-via-go

func makeSigner(keyname string) (signer ssh.Signer, err error) {
	fp, err := os.Open(keyname)
	if err != nil {
		return
	}
	defer fp.Close()

	buf, _ := ioutil.ReadAll(fp)
	signer, _ = ssh.ParsePrivateKey(buf)
	return
}

func authKey() (auth ssh.AuthMethod) {
        signers := []ssh.Signer{}
	keys := []string{os.Getenv("HOME") + "/.ssh/id_rsa", os.Getenv("HOME") + "/.ssh/id_dsa"}

	for _, keyname := range keys {
		signer, err := makeSigner(keyname)
		if err == nil {
			signers = append(signers, signer)
		}
	}

	return ssh.PublicKeys(signers...)
}

func Config(username string, unixConn net.Conn) *ssh.ClientConfig {
	config := &ssh.ClientConfig{
		User: username,
	}
	if os.Getenv("SSH_AUTH_SOCK") != "" {	 			      
		config.Auth = append(config.Auth, authAgent(unixConn))
	}
	config.Auth = append(config.Auth, authKey())

	return config
}
