/*
MIT License

Copyright (c) 2016 FoxBoxsnet

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/
package main

import (
    "os"
    "strings"
    "log"
    "time"

    "github.com/comail/colog"
    "github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-sdk-go/aws/session"
    "github.com/aws/aws-sdk-go/service/route53"
)

var (
    domain          string
    HostedZone      string
    TXT_CHALLENGE   string
)

const (
    sleeptime  = 10
    maxRetries = 5
    route53TTL = 10
)

func main() {
    colog.Register()
    if len(os.Args) <= 4 {
        log.Println("ERROR: parameter")
        log.Println("USAGE:https://github.com/FoxBoxsnet/letsencrypt.sh-dns-route53")
        log.Println("      $./letsencrypt.sh --cron --domain example.com --challenge dns-01 --hook ./letsencrypt.sh-dns-route53")
        os.Exit(1)
    }

        if "deploy_challenge" == os.Args[1] {
            log.Printf("INFO: deploy_challenge common_name: [%s]", os.Args[2])
            common_name := os.Args[2]
            domain := Getdomain(common_name)
            HostedZone := ListHostedZonesByName(domain)
            TXT_CHALLENGE := os.Args[4]
            ChangeResourceRecordSets("UPSERT", common_name, TXT_CHALLENGE, HostedZone, route53TTL)
            os.Exit(0)
        } else if "clean_challenge" == os.Args[1] {
            log.Printf("INFO: clean_challenge common_name: [%s]", os.Args[2])
            common_name := os.Args[2]
            domain := Getdomain(common_name)
            HostedZone := ListHostedZonesByName(domain)
            TXT_CHALLENGE := os.Args[4]
            ChangeResourceRecordSets("DELETE", common_name, TXT_CHALLENGE, HostedZone, route53TTL)
            os.Exit(0)
        } else if "deploy_cert" == os.Args[1] {
            log.Printf("INFO: deploy_cert common_name: [%s]", os.Args[2])
            log.Printf("INFO: Private key       : %v", string(os.Args[3]))
            log.Printf("INFO: Private cert      : %v", string(os.Args[4]))
            log.Printf("INFO: Private fullchain : %v", string(os.Args[5]))
            log.Printf("INFO: Private chain     : %v", string(os.Args[6]))
            os.Exit(0)
        } else if "unchanged_cert" == os.Args[1] {
            log.Printf("INFO: unchanged cert common_name: [%s]", os.Args[2])
            os.Exit(0)
        } else {
            log.Println("ERROR: parameter")
            log.Println("USAGE:https://github.com/FoxBoxsnet/letsencrypt.sh-dns-route53")
            log.Println("      $./letsencrypt.sh --cron --domain example.com --challenge dns-01 --hook ./letsencrypt.sh-dns-route53")
            os.Exit(1)
        }
}

func Getdomain(fqdn string)  string  {
    var domain_tmp1 = []string{}
    domain_tmp1   = strings.Split(fqdn , ".")
    domain_tmp2  := domain_tmp1[len(domain_tmp1)-2:]
    domain       := string(strings.Join(domain_tmp2 , ".") + ".")
    return domain
}

func ListHostedZonesByName(domain string) string {
    sess, err := session.NewSession()
    if err != nil {
        log.Printf("ERORR: failed to create session,", err)
        os.Exit(1)
    }

    svc := route53.New(sess)

    ListHZBN_params := &route53.ListHostedZonesByNameInput{
        DNSName:      aws.String(domain),
    }
    resp, err := svc.ListHostedZonesByName(ListHZBN_params)
    if err != nil {
        log.Printf("ERORR: failed to ListHostedZonesByName,", err)
        os.Exit(1)
    }

    var hostedZoneID string
    for _, hostedZone := range resp.HostedZones {
        if !*hostedZone.Config.PrivateZone && *hostedZone.Name == domain {
            hostedZoneID = *hostedZone.Id
            break
        }
    }

    log.Printf("INFO: Found ZONE : %s Domain : %s\n",hostedZoneID,domain)
    return hostedZoneID
}

func GetChange(GetChangeID string) {
    log.Println("INFO: apply wait please wait for a while...")
    sess, err := session.NewSession()
    if err != nil {
        log.Printf("ERORR: failed to create session,", err)
        os.Exit(1)
    }
    svc := route53.New(sess)

    params := &route53.GetChangeInput{
        Id: aws.String(GetChangeID), // Required
    }

    for i := 0; i < maxRetries; i++ {
        time.Sleep(sleeptime * time.Second)
        resp, err := svc.GetChange(params)
        if err != nil {
            log.Printf("ERORR: failed to create session,", err)
            os.Exit(1)
        }

        if "INSYNC" == *resp.ChangeInfo.Status {
            var status string
            status = *resp.ChangeInfo.Status
            log.Printf("INFO; Record status : %v", status)
            break
        }
    }
}

func ChangeResourceRecordSets (action ,domain , txt_challenge , HostedZones string , ttl int64)  {
    sess, err := session.NewSession()
    if err != nil {
        log.Printf("ERORR: failed to create session,", err)
        os.Exit(1)
    }
    svc := route53.New(sess)

    recordSet := ResourceRecordSet(domain, txt_challenge , ttl)
    reqParams := &route53.ChangeResourceRecordSetsInput{
        ChangeBatch: &route53.ChangeBatch{
            Comment: aws.String("Managed by letsencrypt.sh-dns-route53"),
            Changes: []*route53.Change{
                {
                    Action:            aws.String(action),
                    ResourceRecordSet: recordSet,
                },
            },
        },
        HostedZoneId: aws.String(HostedZones),
    }

    ResourceRecordResp, err := svc.ChangeResourceRecordSets(reqParams)
    if err != nil {
        log.Printf("ERORR: failed to create session,", err)
        os.Exit(1)
    }

    var getChange string
    getChange = *ResourceRecordResp.ChangeInfo.Id
    GetChange(getChange)
}

func ResourceRecordSet(domain , txt_challenge string, ttl int64) *route53.ResourceRecordSet {
    return &route53.ResourceRecordSet{
        Name: aws.String("_acme-challenge." + domain),
        Type: aws.String("TXT"),
        TTL:  aws.Int64(int64(ttl)),
        ResourceRecords: []*route53.ResourceRecord{
            {
                Value: aws.String(`"` + txt_challenge + `"`),
            },
        },
    }
}