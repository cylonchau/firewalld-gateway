(window["webpackJsonp"]=window["webpackJsonp"]||[]).push([["chunk-44e2ecd4"],{"0000":function(e,t,i){"use strict";i.r(t);var n=function(){var e=this,t=e.$createElement,i=e._self._c||t;return i("div",{staticClass:"components-container"},[e._m(0),i("el-form",{ref:"dataForm",staticStyle:{width:"400px","margin-left":"50px"},attrs:{rules:e.rules,model:e.hostRule,"label-position":"left","label-width":"85px"}},[i("el-form-item",{attrs:{label:"IP Range",prop:"ip_range"}},[i("el-input",{model:{value:e.hostRule.ip_range,callback:function(t){e.$set(e.hostRule,"ip_range",t)},expression:"hostRule.ip_range"}})],1),i("el-form-item",{attrs:{label:"Tag",prop:"tag_id"}},[i("el-select",{staticClass:"filter-item",attrs:{placeholder:"Please select"},model:{value:e.hostRule.tag_id,callback:function(t){e.$set(e.hostRule,"tag_id",t)},expression:"hostRule.tag_id"}},e._l(e.calendarTypeOptions,(function(e){return i("el-option",{key:e.id,attrs:{label:e.name,value:e.id}})})),1)],1),i("el-form-item",[i("el-button",{attrs:{type:"primary"},on:{click:function(t){return e.createData()}}},[e._v("Create")]),i("el-button",{on:{click:e.resetForm}},[e._v("Reset")])],1)],1)],1)},a=[function(){var e=this,t=e.$createElement,i=e._self._c||t;return i("aside",[e._v(" Uranus offers automatic network discovery functionality that is effective and very flexible. "),i("p",[e._v("With network discovery properly set up you can:")]),i("ul",[i("li",[e._v("speed up host create.")]),i("li",[e._v("simplify administration")])]),i("span",[e._v("Tips：Network discovery is once task, active host will add in Uranus")])])}],r=i("b775"),s=i("0264"),o={name:"HostDiscovery",data:function(){return{content:'<h1 style="text-align: center;">Welcome to the TinyMCE demo!</h1><p style="text-align: center; font-size: 15px;"><img title="TinyMCE Logo" src="//www.tinymce.com/images/glyph-tinymce@2x.png" alt="TinyMCE Logo" width="110" height="97" /><ul>\n            <li>Our <a href="//www.tinymce.com/docs/">documentation</a> is a great resource for learning how to configure TinyMCE.</li><li>Have a specific question? Visit the <a href="https://community.tinymce.com/forum/">Community Forum</a>.</li><li>We also offer enterprise grade support as part of <a href="https://tinymce.com/pricing">TinyMCE premium subscriptions</a>.</li>\n        </ul>',listLoading:!0,calendarTypeOptions:[],hostRule:{ip_range:"",tag_id:void 0},rules:{ip_range:[{required:!0,validator:function(e,t,i){Object(s["b"])(t)?i():i(t)},message:"Illegal CIDR",trigger:"blur"}],tag_id:[{required:!0,message:"tag_id is required",trigger:"change"}]}}},created:function(){this.getTags()},methods:{getTags:function(){var e=this;this.listLoading=!0,Object(r["g"])("/fw/tag",{offset:0,limit:9999}).then((function(t){e.calendarTypeOptions=t.data.list,e.listLoading=!1}))},createData:function(){var e=this;this.$refs["dataForm"].validate((function(t){t&&Object(r["c"])("/fw/host/async",e.hostRule).then((function(t){e.$notify({title:"Success",message:"Created Successfully",type:"success",duration:2e3})}))}))},resetForm:function(){this.$refs["dataForm"].resetFields()}}},l=o,c=(i("9baf"),i("2877")),u=Object(c["a"])(l,n,a,!1,null,"462b8b77",null);t["default"]=u.exports},"0264":function(e,t,i){"use strict";i.d(t,"a",(function(){return n})),i.d(t,"b",(function(){return a}));i("ac1f"),i("00b4");function n(e){var t,i=new Array;return i[0]=e>>>24>>>0,i[1]=e<<8>>>24>>>0,i[2]=e<<16>>>24,i[3]=e<<24>>>24,t=String(i[0])+"."+String(i[1])+"."+String(i[2])+"."+String(i[3]),t}var a=function(e){return!!/^((\d|[1-9]\d|1\d\d|2([0-4]\d|5[0-5]))\.){3}((\d|[1-9]\d|1\d\d|2([0-4]\d|5[0-5]))\/){1}(\d|[1-3]\d)$/.test(e)}},"9baf":function(e,t,i){"use strict";i("eac3")},eac3:function(e,t,i){}}]);