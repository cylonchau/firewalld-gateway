package apis

import (
	"errors"
	"reflect"
	"strings"
)

type ServiceSetting struct {
	Version      string            `form:"version" json:"version,omitempty"`
	Short        string            `form:"short" json:"short,omitempty"`
	Description  string            `form:"description" json:"description,omitempty"`
	Port         *[]Port           `form:"port" json:"port,omitempty"`
	Module       []string          `form:"module" json:"module,omitempty"`
	Destination  map[string]string `form:"destination" json:"destination,omitempty"`
	Protocol     []string          `form:"protocol" json:"protocol,omitempty"`
	Source_ports []string          `form:"source_ports" json:"source_ports,omitempty"`
}

type Source struct {
	Address string `form:"address" json:"address,omitempty"`
	Mac     string `form:"mac" json:"mac,omitempty"`
	Ipset   string `form:"ipset" json:"ipset,omitempty"`
	Invert  string `form:"invert" json:"invert,omitempty"`
}

type Destination struct {
	Address string `form:"address" json:"address,omitempty"`
	Invert  string `form:"invert" json:"invert,omitempty"`
}

type Port struct {
	Port     string `form:"port" query:"port" json:"port,omitempty"`
	Protocol string `form:"protocol" query:"protocol" json:"protocol,omitempty"`
}
type Protocol struct {
	Value string `form:"value" json:"value,omitempty"`
}

type IcmpBlock struct {
	Name string `form:"name" json:"name,omitempty"`
}
type IcmpType struct {
	Name string `form:"name" json:"name,omitempty"`
}

type ForwardPort struct {
	Port     string `form:"port" json:"port,omitempty"`
	Protocol string `form:"protocol" json:"protocol,omitempty"`
	ToPort   string `form:"toport" json:"toport,omitempty"`
	ToAddr   string `form:"toaddr" json:"toaddr,omitempty"`
}

type Log struct {
	Prefix string `form:"prefix" json:"prefix,omitempty"`
	Level  string `form:"level" json:"level,omitempty"`
	Limit  Limit  `form:"limit" json:"limit,omitempty"`
}
type Limit struct {
	Value string `form:"value" json:"value,omitempty"`
}
type Audit struct {
	Limit Limit `form:"limit" json:"limit,omitempty"`
}
type Accept struct {
	Flag  bool
	Limit Limit `form:"limit" json:"limit,omitempty"`
}
type Reject struct {
	Type  string `form:"type" json:"type,omitempty"`
	Limit Limit  `form:"limit" json:"limit,omitempty"`
}
type Drop struct {
	Flag  bool
	Limit Limit `form:"limit" json:"limit,omitempty"`
}

type Mark struct {
	Set   string `form:"set" json:"set,omitempty"`
	Limit Limit  `form:"limit" json:"limit,omitempty"`
}

type SourcePort struct {
	Port     string `form:"port" json:"port,omitempty"`
	Protocol string `form:"protocol" json:"protocol,omitempty"`
}

type Rule struct {
	Family      string       `form:"family" json:"family,omitempty,default=ipv4"`
	Source      *Source      `form:"source" json:"source,omitempty"`
	Destination *Destination `form:"destination" json:"destination,omitempty"`
	Service     []string     `form:"service" json:"service,omitempty"`
	Port        *Port        `form:"protocol" json:"protocol,omitempty"`
	Protocol    *Protocol    `form:"protocol" json:"protocol,omitempty"`
	IcmpBlock   *IcmpBlock   `form:"icmpblock" json:"icmpblock,omitempty"`
	IcmpType    *IcmpType    `form:"icmptype" json:"icmptype,omitempty"`
	ForwardPort *ForwardPort `form:"forwardport" json:"forwardport,omitempty"`
	Log         *Log         `form:"log" json:"log,omitempty"`
	Audit       *Audit       `form:"audit" json:"audit,omitempty"`
	Accept      *Accept      `form:"accept" json:"accept,omitempty"`
	Reject      *Reject      `form:"reject" json:"reject,omitempty"`
	Drop        *Drop        `form:"drop" json:"drop,omitempty"`
	Mark        *Mark        `form:"mark" json:"mark,omitempty"`
}

type Interface struct {
	Name string `form:"source" json:"name,omitempty"`
}

/*
 * 对应firewalld zoneSettingd的顺序
   [
	   "", version
	   "", short
	   "", description
	   False,  Forward
	   DEFAULT_ZONE_TARGET,  target
	   [], service
	   [], port
	   [], icmp-blocks
	   False,  masquerade
	   [], forward-ports
	   [], interface
	   [], sources
	   [], rich
	   [], protocols
	   [], source-ports
	   False icmp-block-inversion
	]

*/

type QuerySettings struct {
	Version            string         `deepcopier:"field:Version" form:"version" json:"version,omitempty"`
	Short              string         `deepcopier:"field:Short" form:"short" json:"short,omitempty" binding:"required"`
	Description        string         `deepcopier:"field:Description" form:"description" json:"description,omitempty" binding:"required"`
	Forward            bool           `deepcopier:"field:Forward" form:"forward" json:"forward,omitempty"`
	Target             string         `deepcopier:"field:Target" form:"target" json:"target,omitempty" binding:"required"`
	Service            []string       `deepcopier:"field:Service" form:"service" json:"service,omitempty"`
	Port               []*Port        `deepcopier:"field:Port" form:"port" json:"port,omitempty"`
	IcmpBlock          []*IcmpBlock   `deepcopier:"field:IcmpBlock" form:"icmpblock" json:"icmpblock,omitempty"`
	Masquerade         bool           `deepcopier:"field:Masquerade" form:"masquerade" json:"masquerade,omitempty"`
	ForwardPort        []*ForwardPort `deepcopier:"field:ForwardPort" form:"forwardport" json:"forwardport,omitempty"`
	Interface          []*Interface   `deepcopier:"field:Interface" form:"interface" json:"interface,omitempty"`
	Source             []*Source      `deepcopier:"field:Source" form:"source" json:"source,omitempty"`
	Rule               []*Rule        `deepcopier:"skip" form:"rule" json:"rule,omitempty"`
	Protocol           []*Protocol    `deepcopier:"field:Protocol" form:"protocol" json:"protocol,omitempty"`
	SourcePort         []*SourcePort  `deepcopier:"field:SourcePort" form:"sourceport" json:"sourceport,omitempty"`
	IcmpBlockInversion bool           `deepcopier:"field:IcmpBlockInversion" form:"icmp-block-inversion" json:"icmp-block-inversion",omitempty"`
}

type Settings struct {
	Version            string         `deepcopier:"field:Version" form:"version" json:"version,omitempty"`
	Short              string         `deepcopier:"field:Short" form:"short" json:"short,omitempty" binding:"required"`
	Description        string         `deepcopier:"field:Description" form:"description" json:"description,omitempty" binding:"required"`
	Forward            bool           `deepcopier:"field:Forward" form:"forward" json:"forward,omitempty"`
	Target             string         `deepcopier:"field:Target" form:"target" json:"target,omitempty" binding:"required"`
	Service            []string       `deepcopier:"field:Service" form:"service" json:"service,omitempty"`
	Port               []*Port        `deepcopier:"field:Port" form:"port" json:"port,omitempty"`
	IcmpBlock          []*IcmpBlock   `deepcopier:"field:IcmpBlock" form:"icmpblock" json:"icmpblock,omitempty"`
	Masquerade         bool           `deepcopier:"field:Masquerade" form:"masquerade" json:"masquerade,omitempty"`
	ForwardPort        []*ForwardPort `deepcopier:"field:ForwardPort" form:"forwardport" json:"forwardport,omitempty"`
	Interface          []*Interface   `deepcopier:"field:Interface" form:"interface" json:"interface,omitempty"`
	Source             []*Source      `deepcopier:"field:Source" form:"source" json:"source,omitempty"`
	Rule               []string       `deepcopier:"skip" form:"rule" json:"rule,omitempty"`
	Protocol           []*Protocol    `deepcopier:"field:Protocol" form:"protocol" json:"protocol,omitempty"`
	SourcePort         []*SourcePort  `deepcopier:"field:SourcePort" form:"sourceport" json:"sourceport,omitempty"`
	IcmpBlockInversion bool           `deepcopier:"field:IcmpBlockInversion" form:"icmp-block-inversion" json:"icmp-block-inversion",omitempty"`
}

func (this *Source) IsEmpty() bool {
	return this == nil
}

func (this *Destination) IsEmpty() bool {
	return this == nil
}

func (this *Port) IsEmpty() bool {
	return this == nil
}

func (this *Protocol) IsEmpty() bool {
	return this == nil
}

func (this *IcmpBlock) IsEmpty() bool {
	return this == nil
}

func (this *IcmpType) IsEmpty() bool {
	return this == nil
}

func (this *Log) IsEmpty() bool {
	return this == nil
}

func (this *ForwardPort) IsEmpty() bool {
	return this == nil
}

func (this *Audit) IsEmpty() bool {
	return this == nil
}

func (this *Accept) IsEmpty() bool {
	return this == nil
}

func (this *Reject) IsEmpty() bool {
	return this == nil
}

func (this *Drop) IsEmpty() bool {
	return this == nil
}

func (this *Mark) IsEmpty() bool {
	return this == nil
}

func (this *Limit) IsEmpty() bool {
	return this == nil
}

func (this *Source) ToString() string {
	var str = " source "
	if this.Address != "" {
		str += "address=" + this.Address
	} else if this.Mac != "" {
		str += "mac=" + this.Mac
	} else {
		str += "ipset=" + this.Ipset
	}
	if this.Invert != "" {
		str += " "
		str += "invert=" + this.Invert
	}
	str += " "
	return str
}

func (this *Destination) ToString() string {
	var str = " destination "
	if this.Address != "" {
		str += "address=" + this.Address
	}
	if this.Invert != "" {
		str += " "
		str += "invert=" + this.Invert
	}
	str += " "
	return str
}

func (this *Port) ToString() string {
	var str = "port "
	if this.Port != "" {
		str += "name=" + this.Port
	}
	if this.Protocol != "" {
		str += "protocol=" + this.Protocol
	}

	str += " "
	return str
}

func (this *Protocol) ToString() string {
	var str = "Protocol "
	if this.Value != "" {
		str += "value=" + this.Value
	}

	str += " "
	return str
}

func (this *IcmpBlock) ToString() string {
	var str = "icmp-block "
	if this.Name != "" {
		str += "name=" + this.Name
	}

	str += " "
	return str
}

func (this *IcmpType) ToString() string {
	var str = "icmp-type "
	if this.Name != "" {
		str += "name=" + this.Name
	}

	str += " "
	return str
}

func (this *ForwardPort) ToString() string {
	var str = "forward-port "

	if this.Port != "" {
		str += "port=" + this.Port
	}

	if this.Protocol != "" {
		str += " "
		str += "protocol=" + this.Protocol
	}

	if this.ToPort != "" {
		str += " "
		str += "to-port=" + this.ToPort
	}

	if this.ToAddr != "" {
		str += " "
		str += "to-addr=" + this.ToAddr
	}

	str += " "
	return str
}

func (this *Log) ToString() string {
	var str = "log"

	if this.Prefix != "" {
		str += " " + "prefix=" + this.Prefix
	}

	if this.Level != "" {
		str += " " + "level=" + this.Level
	}

	if !this.Limit.IsEmpty() {
		str += " " + "limit value=" + this.Limit.Value
	}

	str += " "
	return str
}

func (this *Audit) ToString() string {
	var str = "audit"

	if !this.Limit.IsEmpty() {
		str += " " + "limit value=" + this.Limit.Value
	}

	str += " "
	return str
}

func (this *Accept) ToString() string {
	var str string

	if this.Flag {
		str = "accept "
	}
	if !this.Limit.IsEmpty() {
		str += "limit value=" + this.Limit.Value
	}

	str += " "
	return str
}

func (this *Reject) ToString() string {
	var str = "reject "

	if this.Type != "" {
		str += "type=" + this.Type
	}

	if !this.Limit.IsEmpty() {
		str += " "
		str += "limit value=" + this.Limit.Value
	}

	str += " "
	return str
}

func (this *Drop) ToString() string {
	var str string

	if this.Flag {
		str = "drop "
	}
	if !this.Limit.IsEmpty() {
		str += "limit value=" + this.Limit.Value
	}
	str += " "
	return str
}

func (this *Mark) ToString() string {
	var str = "mark"

	if this.Set != "" {
		str += " "
		str += "set=" + this.Set
	}

	if !this.Limit.IsEmpty() {
		str += " "
		str += "limit value=" + this.Limit.Value
	}

	str += " "
	return str
}

func (this *Rule) ToString() (ruleString string) {
	ruleString = "rule "
	if this.Family != "" {
		ruleString += "family=" + this.Family
	}

	if !this.Source.IsEmpty() {
		ruleString += this.Source.ToString()
	}

	if !this.Destination.IsEmpty() {
		ruleString += this.Destination.ToString()
	}

	if len(this.Service) > 0 {
		ruleString += "service name=" + this.Service[0] + " "
	}

	if !this.Port.IsEmpty() {
		ruleString += this.Port.ToString()
	}

	if !this.Protocol.IsEmpty() {
		ruleString += this.Protocol.ToString()
	}

	if !this.IcmpBlock.IsEmpty() {
		ruleString += this.IcmpBlock.ToString()
	}

	if !this.IcmpType.IsEmpty() {
		ruleString += this.IcmpType.ToString()
	}

	if !this.ForwardPort.IsEmpty() {
		ruleString += this.ForwardPort.ToString()
	}

	if !this.Log.IsEmpty() {
		ruleString += this.Log.ToString()
	}

	if !this.Audit.IsEmpty() {
		ruleString += this.Audit.ToString()
	}

	if !this.Accept.IsEmpty() {
		ruleString += this.Accept.ToString()
	}

	if !this.Reject.IsEmpty() {
		ruleString += this.Reject.ToString()
	}

	if !this.Drop.IsEmpty() {
		ruleString += this.Drop.ToString()
	}

	if !this.Mark.IsEmpty() {
		ruleString += this.Mark.ToString()
	}
	return
}

func SliceToStruct(array interface{}) (forwardPort ForwardPort, err error) {
	forwardPort = ForwardPort{}
	valueOf := reflect.ValueOf(forwardPort)
	if valueOf.Kind() != reflect.Ptr {
		return ForwardPort{}, errors.New("must ptr")
	}
	valueOf = valueOf.Elem()
	if valueOf.Kind() != reflect.Struct {
		return ForwardPort{}, errors.New("must struct")
	}

	switch array.(type) {
	case []string:
		arrayImplement := array.([]string)
		for i := 0; i < valueOf.NumField(); i++ {
			if i >= len(arrayImplement) {
				break
			}
			val := arrayImplement[i]
			if val != "" && reflect.ValueOf(val).Kind() == valueOf.Field(i).Kind() {
				valueOf.Field(i).Set(reflect.ValueOf(val))
			}
		}
	case []interface{}:
		arrayImplement := array.([]interface{})
		for i := 0; i < valueOf.NumField(); i++ {
			if i >= len(arrayImplement) {
				break
			}
			val := arrayImplement[i]
			if val != "" && reflect.ValueOf(val).Kind() == valueOf.Field(i).Kind() {
				valueOf.Field(i).Set(reflect.ValueOf(val))
			}
		}
	}

	return forwardPort, nil
}

func stringToReject(slice []string) (reject *Reject, ruleSlice []string) {
Label:
	for index, value := range slice {
		tmp_slice := strings.Split(value, "=")
		switch tmp_slice[1] {
		case "type":
			slice = removeSliceElement(slice, index)
			reject.Type = slice[index]
			slice = removeSliceElement(slice, index)
			goto Label
		case "limit":
			slice = removeSliceElement(slice, index)
			tmp_slice := strings.Split(slice[index], "=")
			reject.Limit = Limit{Value: tmp_slice[1]}
			slice = removeSliceElement(slice, index)
			goto Label
		}
	}
	ruleSlice = slice
	return reject, ruleSlice
}

func stringToMark(slice []string) (mark *Mark, ruleSlice []string) {

Label:
	for index, value := range slice {
		tmp_slice := strings.Split(value, "=")
		switch tmp_slice[0] {
		case "set":
			slice = removeSliceElement(slice, index)
			mark.Set = tmp_slice[1]
			goto Label
		case "limit":
			slice = removeSliceElement(slice, index)
			tmp_slice := strings.Split(slice[index], "=")
			mark.Limit = Limit{Value: tmp_slice[1]}
			slice = removeSliceElement(slice, index)
			goto Label
		}
	}
	ruleSlice = slice
	return
}

func stringToForwardPort(slice []string) (forwardPort *ForwardPort, ruleSlice []string) {

Label:
	for index, value := range slice {
		tmp_slice := strings.Split(value, "=")
		switch tmp_slice[0] {
		case "port":
			slice = removeSliceElement(slice, index)
			forwardPort.Port = tmp_slice[1]
			goto Label
		case "protocol":
			slice = removeSliceElement(slice, index)
			forwardPort.Protocol = tmp_slice[1]
			goto Label
		case "to-port":
			slice = removeSliceElement(slice, index)
			forwardPort.ToPort = tmp_slice[1]
			goto Label
		case "to-addr":
			slice = removeSliceElement(slice, index)
			forwardPort.ToAddr = tmp_slice[1]
			goto Label
		}
	}
	ruleSlice = slice
	return
}

func stringToLog(slice []string) (log *Log, ruleSlice []string) {

Label:
	for index, value := range slice {
		tmp_slice := strings.Split(value, "=")
		switch tmp_slice[0] {
		case "prefix":
			slice = removeSliceElement(slice, index)
			log.Prefix = tmp_slice[1]
			goto Label
		case "level":
			slice = removeSliceElement(slice, index)
			log.Level = tmp_slice[1]
			goto Label
		case "limit":
			slice = removeSliceElement(slice, index)
			tmp_slice := strings.Split(slice[index], "=")
			log.Limit = Limit{Value: tmp_slice[1]}
			slice = removeSliceElement(slice, index)
			goto Label
		}
	}
	ruleSlice = slice
	return
}

func StringToRule(str string) (rule *Rule) {

	strslice := strings.Split(str, " ")
	rule = &Rule{}
Label:
	for index, value := range strslice {
		switch value {
		case "rule":
			strslice = removeSliceElement(strslice, index)
			tmp_str := strings.Split(strslice[index], "=")
			rule.Family = tmp_str[1]
			strslice = removeSliceElement(strslice, index)
			goto Label
		case "source":
			strslice = removeSliceElement(strslice, index)
			tmp_str := strings.Split(strslice[index], "=")
			source := &Source{}
			switch tmp_str[0] {
			case "address":
				source.Address = tmp_str[1]
			case "mac":
				source.Mac = tmp_str[1]
			case "ipset":
				source.Ipset = tmp_str[1]
			case "invert":
				source.Invert = tmp_str[1]
			}
			rule.Source = source
			strslice = removeSliceElement(strslice, index)
			goto Label
		case "destination":
			strslice = removeSliceElement(strslice, index)
			dst := &Destination{}
			tmp_str := strings.Split(strslice[index], "=")
			switch tmp_str[0] {
			case "address":
				dst.Address = tmp_str[1]
			case "invert":
				dst.Invert = tmp_str[1]
			}
			rule.Destination = dst
			strslice = removeSliceElement(strslice, index)
			goto Label
		case "service":
			strslice = removeSliceElement(strslice, index)
			tmp_str := strings.Split(strslice[index], "=")
			rule.Service = []string{tmp_str[1]}
			strslice = removeSliceElement(strslice, index)
			goto Label
		case "port":
			strslice = removeSliceElement(strslice, index)
			port := strings.Split(strslice[index], "=")
			protocol := strings.Split(strslice[index+1], "=")
			rule.Port = &Port{
				Port:     port[1],
				Protocol: protocol[1],
			}
			strslice = removeSliceElement(strslice, index)
			goto Label
		case "protocol":
			strslice = removeSliceElement(strslice, index)
			tmp_str := strings.Split(strslice[index+1], "=")
			rule.Protocol = &Protocol{Value: tmp_str[1]}
			strslice = removeSliceElement(strslice, index)
			goto Label
		case "icmp-block":
			strslice = removeSliceElement(strslice, index)
			tmp_str := strings.Split(strslice[index+1], "=")
			rule.IcmpBlock = &IcmpBlock{Name: tmp_str[1]}
			strslice = removeSliceElement(strslice, index)
			goto Label
		case "icmp-type":
			strslice = removeSliceElement(strslice, index)
			tmp_str := strings.Split(strslice[index+1], "=")
			rule.IcmpType = &IcmpType{Name: tmp_str[1]}
			strslice = removeSliceElement(strslice, index)
			goto Label
		case "forward-port":
			strslice = removeSliceElement(strslice, index)
			rule.ForwardPort, strslice = stringToForwardPort(strslice)
			goto Label
		case "log":
			strslice = removeSliceElement(strslice, index)
			rule.Log, strslice = stringToLog(strslice)
			goto Label
		case "audit":
			strslice = removeSliceElement(strslice, index)
			tmp_str := strings.Split(strslice[index], "=")
			rule.Audit = &Audit{Limit: Limit{Value: tmp_str[1]}}
			strslice = removeSliceElement(strslice, index)
			goto Label
		case "accept":
			strslice = removeSliceElement(strslice, index)
			accept := &Accept{}
			rule.Accept = accept
			rule.Accept.Flag = true
			var tmp_str []string
			if len(strslice) > 0 {
				if strslice[index] == "limit" {
					strslice = removeSliceElement(strslice, index)
					tmp_str = strings.Split(strslice[index], "=")
					rule.Accept.Limit = Limit{Value: tmp_str[1]}
				}
			}
			goto Label
		case "drop":
			var tmp_str []string
			strslice = removeSliceElement(strslice, index)
			rule.Drop.Flag = true
			if len(strslice) > 0 {
				if strslice[index] == "limit" {
					strslice = removeSliceElement(strslice, index)
					tmp_str = strings.Split(strslice[index], "=")
					rule.Drop.Limit = Limit{Value: tmp_str[1]}
				}
			}
			goto Label
		case "reject":
			strslice = removeSliceElement(strslice, index)
			rule.Reject, strslice = stringToReject(strslice)
			goto Label
		case "mark":
			strslice = removeSliceElement(strslice, index)
			rule.Mark, strslice = stringToMark(strslice)
			goto Label
		}
	}
	return rule
}

func removeSliceElement(s []string, index int) []string {
	return append(s[:index], s[index+1:]...)
}
