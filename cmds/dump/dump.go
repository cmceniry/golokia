package main

import (
	"fmt"
	"github.com/swarvanusg/golokia"
	"os"
)

func main() {
	if len(os.Args) != 6 {
		fmt.Fprintf(os.Stderr, "Invalid command line : must specify jolokia target (host, port, jolokiaServiceName, targetHost, targetPort)")
		os.Exit(-1)
	}
	jc := golokia.NewJolokiaClient("http://" + os.Args[1] + ":" + os.Args[2] + "/" + os.Args[3])
	jc.SetTarget(os.Args[4] + ":" + os.Args[5])
	domains, err := jc.ListDomains()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to get Domains : %s", err)
		os.Exit(-2)
	}
	for _, d := range domains {
		beans, err := jc.ListBeans(d)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Unable to get beans for %s Domain : %s\n", d, err)
			continue
		}
		for _, b := range beans {
			props, err := jc.ListProperties(d, []string{b})
			if err != nil {
				fmt.Fprintf(os.Stderr, "Unable to get properties for %s,%s Bean : %s\n", d, b, err)
				continue
			}
			for _, p := range props {
				val, err := jc.GetAttr(d, []string{b}, p)
				if err != nil {
					fmt.Fprintf(os.Stderr, "Unable to get value for %s,%s,%s Property : %s\n", d, b, p, err)
					continue
				}
				fmt.Printf("%s,%s,%s = %v\n", d, b, p, val)
			}
		}
	}
}
