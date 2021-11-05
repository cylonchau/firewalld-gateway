package zone

import (
	"firewalld/type"
	"fmt"
	"github.com/godbus/dbus/v5"
	"github.com/joncalhoun/qson"
	"log"
	"net"
	"reflect"
	"regexp"
	"strings"
	"sync"
)

func init() {
	log.SetFlags(log.Ldate | log.Lshortfile)
}

type (
	Rish struct {
		Conn *dbus.Conn
	}

	Rich struct {
		Address  string `form:"address" json:"address" binding:"required,omitempty"`
		Port     int    `form:"port" json:"port" binding:"required,omitempty"`
		Protocol string `form:"protocol" json:"protocol" binding:"required,omitempty"`
		Type     string `form:"type" json:"type" binding:"required,omitempty"`
		Expire   int    `form:"expire" json:"expire,omitempty"`
		Zone     string `form:"utils" json:"utils,omitempty"`
	}

	QueryRich struct {
		Address  string `form:"address" json:"address" binding:"required,omitempty"`
		Port     int    `form:"port" json:"port" binding:"required,omitempty"`
		Protocol string `form:"protocol" json:"protocol" binding:"required,omitempty"`
		Type     string `form:"type" json:"type" binding:"required,omitempty"`
		Expire   int    `form:"expire" json:"expire,omitempty"`
		Zone     string `form:"utils" json:"utils" binding:"required,omitempty"`
	}
)

var (
	cli *dbus.Conn
	err error
)

var instance *Rish
var once sync.Once

func GetInstance() *Rish {
	once.Do(func() {
		instance = &Rish{}
		instance.Conn, err = dbus.SystemBus()
		if err != nil {
			panic(err)
		}
	})
	return instance
}

func (this *Rish) structToRich(r Rich) string {
	var (
		family  = "ipv4"
		rulestr string
	)

	ip, _, _ := net.ParseCIDR(r.Address)
	_, err := net.ParseMAC(r.Address)

	if ip == nil {

	} else if net.IP.To4(ip) == nil {
		family = "ipv6"
	}

	if err == nil {
		re := regexp.MustCompile(`\-|\.`)
		r.Address = re.ReplaceAllString(r.Address, ":")
		rulestr = fmt.Sprintf("rule family=%s source mac=%s port port=%d protocol=%s %s", family, r.Address, r.Port, r.Protocol, r.Type)
	} else {
		rulestr = fmt.Sprintf("rule family=%s source address=%s port port=%d protocol=%s %s", family, r.Address, r.Port, r.Protocol, r.Type)
	}

	return rulestr
}

func (this *Rish) structToQueryRich(r QueryRich) string {
	var (
		family  = "ipv4"
		richstr string
	)

	ip, _, _ := net.ParseCIDR(r.Address)
	_, err := net.ParseMAC(r.Address)

	if ip == nil {

	} else if net.IP.To4(ip) == nil {
		family = "ipv6"
	}

	if err == nil {
		re := regexp.MustCompile(`\-|\.`)
		r.Address = re.ReplaceAllString(r.Address, ":")
		richstr = fmt.Sprintf("rule family=%s source mac=%s port port=%d protocol=%s %s", family, r.Address, r.Port, r.Protocol, r.Type)
	} else {
		richstr = fmt.Sprintf("rule family=%s source address=%s port port=%d protocol=%s %s", family, r.Address, r.Port, r.Protocol, r.Type)
	}

	return richstr
}

/*
 *  insert rich list into rule.
 *  @return error
 */

func (this *Rish) AddRich(rule Rich) (err error) {
	if rule.Zone == "" {
		rule.Zone = "public"
	}

	obj := Conn.Object(_type.Interfaces, _type.Path)

	rulestr := structToRich(rule)
	callReply := obj.Call(_type.AddRich, dbus.FlagNoAutoStart, rule.Zone, rulestr, rule.Expire)

	if callReply.Err != nil {
		err = callReply.Err
		return
	}
	return
}

/*
 *  quert whether rich exists.
 *  @return bool
 */

func (this *Rish) richExists(rule QueryRich) bool {
	obj := Conn.Object(_type.Interfaces, _type.Path)
	richstr := structToQueryRich(rule)
	callReply := obj.Call(_type.QueryRich, dbus.FlagNoAutoStart, rule.Zone, richstr)

	if callReply.Err != nil {
		log.Print("query result: ", err)
		return false
	}

	return true
}

/*
 *  quert whether rich exists.
 *  @return bool
 */

func (this *Rish) QueryRich(rule QueryRich) (isExists bool) {
	richstr := structToQueryRich(rule)
	obj := Conn.Object(_type.Interfaces, _type.Path)
	callReply := obj.Call(_type.QueryRich, dbus.FlagNoAutoStart, rule.Zone, richstr)

	if callReply.Err != nil {
		log.Print("query result: ", err)
		return false
	}

	for _, v := range callReply.Body {
		str := reflect.ValueOf(v)
		isExists = str.Interface().(bool)
	}

	return
}

/*
 *  Converts dbus reply content interface{} to json.
 *  @return []string json string
 */

func (this *Rish) listRichParse(lists []string) (jsonList []string, err error) {

	var (
		re    = regexp.MustCompile(`port\s|source\s|rule\s|service\s|\"`)
		reuri = regexp.MustCompile(`\s`)
	)

	for _, v := range lists {
		v = re.ReplaceAllString(v, "")
		var s string
		if strings.Contains(v, "accept") {
			s = strings.Replace(v, "accept", "type=accept", -1)
		} else if strings.Contains(s, "reject") {
			s = strings.Replace(s, "reject", "type=reject", -1)
		}
		s = reuri.ReplaceAllString(s, "&")

		b, err := qson.ToJSON(s)

		if err != nil {
			log.Print("rich rule conver to string fail: ", err)
			return nil, err
		}

		jsonList = append(jsonList, string(b))
	}

	return
}

/*
 *  Converts dbus reply content interface{} to json.
 *  @return string json []byte
 */

func (this *Rish) oneRichParse(richs []string) (jsonrich []byte) {

	var (
		re    = regexp.MustCompile(`port\s|source\s|rule\s|service\s|\,"`)
		reuri = regexp.MustCompile(`\s`)
	)

	for _, v := range richs {
		v = re.ReplaceAllString(v, "")
		var s string
		if strings.Contains(v, "accept") {
			s = strings.Replace(v, "accept", "type=accept", -1)
		} else if strings.Contains(s, "reject") {
			s = strings.Replace(s, "reject", "type=reject", -1)
		}
		s = reuri.ReplaceAllString(s, "&")

		jsonrich, err := qson.ToJSON(s)
		if err != nil {
			log.Print("rich rule conver to string fail: ", err)
			return nil
		}
		jsonrich = jsonrich
	}

	return jsonrich
}

/*
 *
 *
 */

func (this *Rish) QueryAllRish() (richs []string, err error) {

	obj := Conn.Object(_type.Interfaces, _type.Path)
	call := obj.Call(_type.GetRichs, dbus.FlagNoAutoStart, "public")

	if call.Err != nil {
		err = call.Err
		return
	}

	var lists []string

	for _, v := range call.Body {
		str := reflect.ValueOf(v)
		lists = str.Interface().([]string)
	}

	return listRichParse(lists)
}

func (this *Rish) DelRich(rule Rich) error {
	obj := Conn.Object(_type.Interfaces, _type.Path)

	if rule.Zone == "" {
		rule.Zone = "public"
	}
	richstr := structToRich(rule)
	call := obj.Call(_type.RemoveRichs, dbus.FlagNoAutoStart, rule.Zone, richstr)
	log.Print("remove rich: ", richstr)
	if call.Err != nil {
		err = call.Err
		return err
	}
	return nil
}
