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
    maxRetries = 5
    route53TTL = 10
)


func main() {
    colog.Register()
    
    domain          := Getdomain(os.Args[2])
    HostedZone      := ListHostedZonesByName(domain)
    TXT_CHALLENGE   := os.Args[4]

    if "deploy_challenge" == os.Args[1] {
        log.Println("INFO: deploy_challenge")
        ChangeResourceRecordSets("UPSERT",domain,TXT_CHALLENGE,HostedZone,route53TTL)

    }else if "clean_challenge" == os.Args[1] {
        log.Println("INFO: clean_challenge")
        ChangeResourceRecordSets("DELETE",domain,TXT_CHALLENGE,HostedZone,route53TTL)

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
        return ""
    }

    svc := route53.New(sess)

    ListHZBN_params := &route53.ListHostedZonesByNameInput{
        DNSName:      aws.String(domain),
    }
    resp, err := svc.ListHostedZonesByName(ListHZBN_params)
    if err != nil {
        log.Printf("ERORR: failed to ListHostedZonesByName,", err)
        return ""
    }

    var hostedZoneID string
    for _, hostedZone := range resp.HostedZones {
        // .Name has a trailing dot
        if !*hostedZone.Config.PrivateZone && *hostedZone.Name == domain {
            hostedZoneID = *hostedZone.Id
            break
        }
    }

    log.Printf("INFO: Found ZONE : %s Domain : %s\n",hostedZoneID,domain)
    return hostedZoneID
}

func GetChange(GetChangeID string) {
    log.Println("INFO: apply wait please wait for a while.")
    sess, err := session.NewSession()
    if err != nil {
        log.Println("failed to create session,", err)
        return
    }
    svc := route53.New(sess)

    params := &route53.GetChangeInput{
        Id: aws.String(GetChangeID), // Required
    }

    for i := 0; i < maxRetries; i++ {
        time.Sleep(10 * time.Second)
        resp, err := svc.GetChange(params)

        if err != nil {
            // Print the error, cast err to awserr.Error to get the Code and
            // Message from an error.
            log.Println(err.Error())
            return
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
        //       return ""
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