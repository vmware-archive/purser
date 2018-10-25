/*
 * Copyright (c) 2018 VMware Inc. All Rights Reserved.
 * SPDX-License-Identifier: Apache-2.0
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *    http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package utils

import (
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/Sirupsen/logrus"
)

// PurgeTCPData handles IP conversion from Hex to Dec and cleans up data to contain only
// inter pod address information.
func PurgeTCPData(data string) []string {
	var tcpDump []string

	tcpDumpHex := getTCPDumpHexFromData(data)
	for _, address := range tcpDumpHex {
		localIP, localPort := hexToDecIP(address[6:14]), address[15:19]
		remoteIP, remotePort := hexToDecIP(address[20:28]), address[29:33]

		if isLocalHost(localIP, remoteIP) {
			continue
		}

		addressMapping := localIP + ":" + localPort + ":" + remoteIP + ":" + remotePort
		tcpDump = append(tcpDump, addressMapping)
	}
	return tcpDump
}

// PurgeTCP6Data handles IP conversion from Hex to Dec and cleans up data to contain only
// inter pod address information.
func PurgeTCP6Data(data string) []string {
	var tcpDump []string

	tcpDumpHex := getTCPDumpHexFromData(data)
	for _, address := range tcpDumpHex {
		localIP, localPort := hexToDecIP(address[30:38]), address[39:43]
		remoteIP, remotePort := hexToDecIP(address[68:76]), address[77:81]

		if isLocalHost(localIP, remoteIP) {
			continue
		}

		addressMapping := localIP + ":" + localPort + ":" + remoteIP + ":" + remotePort
		tcpDump = append(tcpDump, addressMapping)
	}
	return tcpDump
}

func getTCPDumpHexFromData(data string) []string {
	tcpDumpHex := strings.Split(data, "\n")
	if len(tcpDumpHex) <= 1 {
		return nil
	}

	// ignore title and last one as it is empty
	tcpDumpHex = tcpDumpHex[1 : len(tcpDumpHex)-1]
	return tcpDumpHex
}

func hexToDecIP(hexIP string) string {
	decBytes, err := hex.DecodeString(hexIP)
	if err != nil {
		logrus.Warnf("failed to decode string to hex %v", err)
	}
	return fmt.Sprintf("%v.%v.%v.%v", decBytes[3], decBytes[2], decBytes[1], decBytes[0])
}

func isLocalHost(localIP, remoteIP string) bool {
	return strings.Compare(localIP, "0.0.0.0") == 0 || strings.Compare(localIP, "127.0.0.1") == 0 || strings.Compare(remoteIP, "0.0.0.0") == 0
}
