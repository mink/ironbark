package main

import (
	"bytes"
	"crypto/tls"
	"encoding/binary"
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"net"
	"time"
)

func main() {
	cert, err := tls.LoadX509KeyPair("cert.pem", "key.pem")
	if err != nil {
		log.Fatalf("failed to load X.509 keypair: %v", err)
	}

	config := &tls.Config{
		Certificates: []tls.Certificate{cert},
	}

	listener, err := tls.Listen("tcp", ":700", config)
	if err != nil {
		panic(err)
	}

	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("failed to accept connection: %v", err)
			continue
		}

		go handle(conn)
	}
}

func handle(conn net.Conn) {
	defer conn.Close()

	tlsConn, ok := conn.(*tls.Conn)
	if ok {
		if err := tlsConn.Handshake(); err != nil {
			fmt.Println("TLS handshake failed:", err)
			return
		}
	}

	// send greeting message

	epp := EPP{
		Xmlns: "urn:ietf:params:xml:ns:epp-1.0",
		Greeting: Greeting{
			SvID:   "Ironbark EPP Server",
			SvDate: time.Now().UTC().Format("2006-01-02T15:04:05.0Z"),
			SvcMenu: ServiceMenu{
				Version: []string{"1.0"},
				Lang:    []string{"en"},
				ObjURI: []string{
					"urn:ietf:params:xml:ns:obj1",
					"urn:ietf:params:xml:ns:obj2",
					"urn:ietf:params:xml:ns:obj3",
				},
				SvcExtension: &SvcExtension{
					ExtURI: []string{"http://custom/obj1ext-1.0"},
				},
			},
			DCP: DCP{
				Access: Access{All: &struct{}{}},
				Statement: []DCPStatement{
					{
						Purpose: DCPPurpose{
							Admin: &struct{}{},
							Prov:  &struct{}{}},
						Recipient: DCPRecipient{
							Ours:   &Ours{},
							Public: &struct{}{},
						},
						Retention: DCPRetention{
							Stated: &struct{}{},
						},
					},
				},
				Expiry: nil,
			},
		},
	}

	marshalled, err := xml.MarshalIndent(epp, "", "  ")
	if err != nil {
		fmt.Printf("Failed to marshal XML: %v", err)
		return
	}
	greeting := string(ConvertSelfClosingTags(marshalled))

	err = writeFramed(conn, greeting)
	if err != nil {
		fmt.Println("Failed to write greeting:", err)
		return
	}

	// receive commands

	for {
		lenBuf := make([]byte, 4)
		if _, err := io.ReadFull(conn, lenBuf); err != nil {
			fmt.Println("Read length error:", err)
			break
		}

		var length int32
		if err := binary.Read(bytes.NewReader(lenBuf), binary.BigEndian, &length); err != nil {
			fmt.Println("Parse length error:", err)
			break
		}

		if length < 4 {
			fmt.Println("Invalid message length:", length)
			break
		}

		bodyLen := length - 4
		buf := make([]byte, bodyLen)
		if _, err := io.ReadFull(conn, buf); err != nil {
			fmt.Println("Read body error:", err)
			break
		}

		fmt.Println("Received:", string(buf))

		// todo: handle commands
		if bytes.Contains(buf, []byte("<login")) {
			break
		}
	}
}

func writeFramed(conn net.Conn, msg string) error {
	buf := new(bytes.Buffer)

	length := int32(4 + len(msg))
	if err := binary.Write(buf, binary.BigEndian, length); err != nil {
		return err
	}

	if _, err := buf.Write([]byte(msg)); err != nil {
		return err
	}

	_, err := conn.Write(buf.Bytes())
	return err
}
