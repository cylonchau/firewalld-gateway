(window["webpackJsonp"]=window["webpackJsonp"]||[]).push([["chunk-677087a4"],{"09f4":function(t,e,a){"use strict";a.d(e,"a",(function(){return r})),Math.easeInOutQuad=function(t,e,a,i){return t/=i/2,t<1?a/2*t*t+e:(t--,-a/2*(t*(t-2)-1)+e)};var i=function(){return window.requestAnimationFrame||window.webkitRequestAnimationFrame||window.mozRequestAnimationFrame||function(t){window.setTimeout(t,1e3/60)}}();function n(t){document.documentElement.scrollTop=t,document.body.parentNode.scrollTop=t,document.body.scrollTop=t}function o(){return document.documentElement.scrollTop||document.body.parentNode.scrollTop||document.body.scrollTop}function r(t,e,a){var r=o(),s=t-r,l=20,c=0;e="undefined"===typeof e?500:e;var u=function t(){c+=l;var o=Math.easeInOutQuad(c,r,s,e);n(o),c<e?i(t):a&&"function"===typeof a&&a()};u()}},2423:function(t,e,a){"use strict";a.d(e,"c",(function(){return n})),a.d(e,"b",(function(){return o})),a.d(e,"d",(function(){return r})),a.d(e,"a",(function(){return s})),a.d(e,"e",(function(){return l}));var i=a("b775");function n(t){return Object(i["f"])({url:"/vue-element-admin/article/list",method:"get",params:t})}function o(t){return Object(i["f"])({url:"/vue-element-admin/article/detail",method:"get",params:{id:t}})}function r(t){return Object(i["f"])({url:"/vue-element-admin/article/pv",method:"get",params:{pv:t}})}function s(t){return Object(i["f"])({url:"/vue-element-admin/article/create",method:"post",data:t})}function l(t){return Object(i["f"])({url:"/vue-element-admin/article/update",method:"post",data:t})}},"35b6":function(t,e,a){},"4e4f":function(t,e,a){"use strict";a("35b6")},"4e82":function(t,e,a){"use strict";var i=a("23e7"),n=a("1c0b"),o=a("7b0b"),r=a("d039"),s=a("a640"),l=[],c=l.sort,u=r((function(){l.sort(void 0)})),d=r((function(){l.sort(null)})),p=s("sort"),f=u||!d||!p;i({target:"Array",proto:!0,forced:f},{sort:function(t){return void 0===t?c.call(o(this)):c.call(o(this),n(t))}})},6724:function(t,e,a){"use strict";a("8d41");var i="@@wavesContext";function n(t,e){function a(a){var i=Object.assign({},e.value),n=Object.assign({ele:t,type:"hit",color:"rgba(0, 0, 0, 0.15)"},i),o=n.ele;if(o){o.style.position="relative",o.style.overflow="hidden";var r=o.getBoundingClientRect(),s=o.querySelector(".waves-ripple");switch(s?s.className="waves-ripple":(s=document.createElement("span"),s.className="waves-ripple",s.style.height=s.style.width=Math.max(r.width,r.height)+"px",o.appendChild(s)),n.type){case"center":s.style.top=r.height/2-s.offsetHeight/2+"px",s.style.left=r.width/2-s.offsetWidth/2+"px";break;default:s.style.top=(a.pageY-r.top-s.offsetHeight/2-document.documentElement.scrollTop||document.body.scrollTop)+"px",s.style.left=(a.pageX-r.left-s.offsetWidth/2-document.documentElement.scrollLeft||document.body.scrollLeft)+"px"}return s.style.backgroundColor=n.color,s.className="waves-ripple z-active",!1}}return t[i]?t[i].removeHandle=a:t[i]={removeHandle:a},a}var o={bind:function(t,e){t.addEventListener("click",n(t,e),!1)},update:function(t,e){t.removeEventListener("click",t[i].removeHandle,!1),t.addEventListener("click",n(t,e),!1)},unbind:function(t){t.removeEventListener("click",t[i].removeHandle,!1),t[i]=null,delete t[i]}},r=function(t){t.directive("waves",o)};window.Vue&&(window.waves=o,Vue.use(r)),o.install=r;e["a"]=o},"8d41":function(t,e,a){},"8fa4":function(t,e,a){"use strict";a.r(e);var i=function(){var t=this,e=t.$createElement,a=t._self._c||e;return a("div",{staticClass:"app-container"},[a("div",{staticClass:"filter-container"},[a("el-input",{staticClass:"filter-item",staticStyle:{width:"200px"},attrs:{placeholder:"Title"},nativeOn:{keyup:function(e){return!e.type.indexOf("key")&&t._k(e.keyCode,"enter",13,e.key,"Enter")?null:t.handleFilter(e)}},model:{value:t.listQuery.title,callback:function(e){t.$set(t.listQuery,"title",e)},expression:"listQuery.title"}}),a("el-button",{directives:[{name:"waves",rawName:"v-waves"}],staticClass:"filter-item",staticStyle:{"margin-left":"10px"},attrs:{type:"primary",icon:"el-icon-search"},on:{click:t.handleFilter}},[t._v(" Search ")]),a("el-button",{staticClass:"filter-item",staticStyle:{"margin-left":"10px"},attrs:{type:"",icon:"el-icon-document-add"},on:{click:t.handleCreate}}),a("el-tag",{directives:[{name:"waves",rawName:"v-waves"}],staticClass:"filter-item",staticStyle:{"margin-left":"10px"},attrs:{size:"35px",effect:"plain",type:""}},[t._v(" NAT Status "),a("el-switch",{staticClass:"switch",attrs:{"active-text":"ON","inactive-text":"OFF","active-value":1,"inactive-value":0},on:{change:t.chageSwitch},model:{value:t.masquerade,callback:function(e){t.masquerade=e},expression:"masquerade"}})],1)],1),a("el-table",{directives:[{name:"loading",rawName:"v-loading",value:t.listLoading,expression:"listLoading"}],key:t.tableKey,staticStyle:{width:"100%"},attrs:{data:t.list.slice((t.currentPage-1)*t.pageSize,t.currentPage*t.pageSize),border:"",fit:"","highlight-current-row":""},on:{"sort-change":t.sortChange}},[a("el-table-column",{attrs:{type:"selection",width:"45"}}),a("el-table-column",{attrs:{label:"Protocol",width:"100px",align:"center"},scopedSlots:t._u([{key:"default",fn:function(e){var i=e.row;return[a("span",[t._v(" "+t._s(i.protocol))])]}}])}),a("el-table-column",{attrs:{label:"Source Port","min-width":"130px",align:"center"},scopedSlots:t._u([{key:"default",fn:function(e){var i=e.row;return[a("el-tag",{attrs:{type:"success",effect:"plain"}},[t._v(" "+t._s(i.port)+" ")])]}}])}),a("el-table-column",{attrs:{label:"Destnation","min-width":"130px",align:"center"},scopedSlots:t._u([{key:"default",fn:function(e){var i=e.row;return[a("el-tag",{attrs:{type:"warning",effect:"plain"}},[t._v(" "+t._s(i.toaddr)+":"+t._s(i.toport)+" ")])]}}])}),a("el-table-column",{attrs:{label:"Actions",align:"center",width:"100","class-name":"small-padding fixed-width"},scopedSlots:t._u([{key:"default",fn:function(e){var i=e.row;return[a("span",[a("el-popconfirm",{attrs:{"confirm-button-text":"Yes","cancel-button-text":"No",icon:"el-icon-info","icon-color":"red",title:"Do you want delete this rule?"},on:{confirm:function(e){return t.handleDelete(i)},onConfirm:function(e){return t.handleDelete(i)}}},[a("el-button",{staticClass:"el-icon-delete",staticStyle:{"margin-left":"5px"},attrs:{slot:"reference",circle:"",type:"danger"},slot:"reference"})],1)],1)]}}])})],1),a("el-pagination",{staticStyle:{"margin-top":"20px"},attrs:{"current-page":t.currentPage,background:"","page-sizes":[5,10,30],"page-size":t.pageSize,layout:"total, sizes,  prev, pager, next, jumper",total:t.list.length},on:{"size-change":t.handleSizeChange,"current-change":t.handleCurrentChange}}),a("el-dialog",{attrs:{title:t.textMap[t.dialogStatus],visible:t.dialogFormVisible},on:{"update:visible":function(e){t.dialogFormVisible=e}}},[a("el-form",{ref:"dataForm",staticStyle:{width:"400px","margin-left":"50px"},attrs:{rules:t.rules,model:t.forward,"label-position":"left","label-width":"115px"}},[a("el-form-item",{attrs:{label:"Protocol",prop:"protocol"}},[a("el-select",{staticClass:"filter-item",attrs:{placeholder:"Please select"},model:{value:t.forward.protocol,callback:function(e){t.$set(t.forward,"protocol",e)},expression:"forward.protocol"}},t._l(t.calendarTypeOptions,(function(t){return a("el-option",{key:t.key,attrs:{label:t.display_name,value:t.key}})})),1)],1),a("el-form-item",{attrs:{label:"Src port",prop:"port"}},[a("el-input",{model:{value:t.forward.port,callback:function(e){t.$set(t.forward,"port",e)},expression:"forward.port"}})],1),a("el-form-item",{attrs:{label:"Dst address",prop:"toaddr"}},[a("el-input",{model:{value:t.forward.toaddr,callback:function(e){t.$set(t.forward,"toaddr",e)},expression:"forward.toaddr"}})],1),a("el-form-item",{attrs:{label:"Dst port",prop:"toport"}},[a("el-input",{model:{value:t.forward.toport,callback:function(e){t.$set(t.forward,"toport",e)},expression:"forward.toport"}})],1),a("el-form-item",{attrs:{label:"Timeout",prop:"timeout"}},[a("el-input",{model:{value:t.temp.timeout,callback:function(e){t.$set(t.temp,"timeout",t._n(e))},expression:"temp.timeout"}})],1)],1),a("div",{staticClass:"dialog-footer",attrs:{slot:"footer"},slot:"footer"},[a("el-button",{on:{click:function(e){t.dialogFormVisible=!1}}},[t._v(" Cancel ")]),a("el-button",{attrs:{type:"primary"},on:{click:function(e){"create"===t.dialogStatus?t.createData():t.updateData()}}},[t._v(" Confirm ")])],1)],1),a("el-dialog",{attrs:{visible:t.dialogPvVisible,title:"Reading statistics"},on:{"update:visible":function(e){t.dialogPvVisible=e}}},[a("el-table",{staticStyle:{width:"100%"},attrs:{data:t.pvData,border:"",fit:"","highlight-current-row":""}},[a("el-table-column",{attrs:{prop:"key",label:"Channel"}}),a("el-table-column",{attrs:{prop:"pv",label:"Pv"}})],1),a("span",{staticClass:"dialog-footer",attrs:{slot:"footer"},slot:"footer"},[a("el-button",{attrs:{type:"primary"},on:{click:function(e){t.dialogPvVisible=!1}}},[t._v("Confirm")])],1)],1)],1)},n=[],o=(a("13d5"),a("d3b7"),a("4e82"),a("c740"),a("a434"),a("d81d"),a("2423")),r=a("6724"),s=a("ed08"),l=a("333d"),c=a("b775"),u=[{key:"tcp",display_name:"TCP"},{key:"udp",display_name:"UDP"}],d=u.reduce((function(t,e){return t[e.key]=e.display_name,t}),{}),p={name:"ComplexTable",components:{Pagination:l["a"]},directives:{waves:r["a"]},filters:{statusFilter:function(t){var e={published:"success",draft:"info",deleted:"danger"};return e[t]},typeFilter:function(t){return d[t]}},data:function(){return{tableKey:0,list:[],total:0,listLoading:!0,listQuery:{page:1,limit:10,importance:void 0,title:void 0,type:void 0,sort:"+id"},masquerade:void 0,importanceOptions:[1,2,3],calendarTypeOptions:u,sortOptions:[{label:"ID Ascending",key:"+id"},{label:"ID Descending",key:"-id"}],statusOptions:["published","draft","deleted"],showReviewer:!1,forward:{port:"",protocol:"",toaddr:"",toport:0},temp:{forward:{port:"",protocol:"",toaddr:"",toport:0},timeout:0,ip:""},dialogFormVisible:!1,dialogStatus:"",textMap:{update:"Edit",create:"Create"},dialogPvVisible:!1,pvData:[],rules:{protocol:[{required:!0,message:"procotol is required",trigger:"change"}],port:[{required:!0,message:"source port is required",trigger:"change"}],toaddr:[{required:!0,message:"destnation address is required",trigger:"change"}],toport:[{required:!0,message:"destnation port is required",trigger:"change"}]},downloadLoading:!1,currentPage:1,pageSize:10}},created:function(){this.getList()},methods:{getMasqueradeStatus:function(){},chageSwitch:function(){var t=this;0==this.$route.params.status?void 0!=this.masquerade&&(0==this.masquerade?Object(c["a"])("/fw/v1/masquerade",{ip:this.$route.params.ip}).then((function(e){t.$notify({title:"Success",message:"Disable Successfully",type:"success",duration:2e3})})):Object(c["d"])("/fw/v1/masquerade",{ip:this.$route.params.ip}).then((function(e){t.$notify({title:"Success",message:"Enable Successfully",type:"success",duration:2e3})}))):0==this.masquerade?Object(c["a"])("/fw/v2/masquerade",{ip:this.$route.params.ip}).then((function(e){t.$notify({title:"Success",message:"Disable Successfully",type:"success",duration:2e3})})):Object(c["d"])("/fw/v2/masquerade",{ip:this.$route.params.ip}).then((function(e){t.$notify({title:"Success",message:"Enable Successfully",type:"success",duration:2e3})}))},handleCurrentChange:function(t){this.currentPage=t},handleSizeChange:function(t){this.pageSize=t},getList:function(){var t=this;this.listLoading=!0,0==this.$route.params.status?(Object(c["g"])("/fw/v1/masquerade",{ip:this.$route.params.ip}).then((function(e){t.masquerade=1==e.data?1:0})),Object(c["g"])("/fw/v1/nat",{ip:this.$route.params.ip}).then((function(e){console.log(e.data),null==e.data?t.list=[]:t.list=e.data}))):(Object(c["g"])("/fw/v2/masquerade",{ip:this.$route.params.ip}).then((function(e){t.masquerade=1==e.data?1:0})),Object(c["g"])("/fw/v2/nat",{ip:this.$route.params.ip}).then((function(e){console.log(e.data),null==e.data?t.list=[]:t.list=e.data}))),this.listLoading=!1},handleFilter:function(){this.listQuery.page=1,this.getList()},handleModifyStatus:function(t,e){this.$message({message:"操作Success",type:"success"}),t.status=e},sortChange:function(t){var e=t.prop,a=t.order;"id"===e&&this.sortByID(a)},sortByID:function(t){this.listQuery.sort="ascending"===t?"+id":"-id",this.handleFilter()},resetTemp:function(){this.forward={port:"",protocol:"",toaddr:"",toport:0}},handleCreate:function(){var t=this;this.resetTemp(),this.dialogStatus="create",this.dialogFormVisible=!0,this.$nextTick((function(){t.$refs["dataForm"].clearValidate()}))},createData:function(){var t=this;this.$refs["dataForm"].validate((function(e){e&&(t.temp.forward=t.forward,t.temp.ip=t.$route.params.ip,0==t.$route.params.status?Object(c["e"])("/fw/v1/nat",t.temp).then((function(e){t.dialogFormVisible=!1,t.getList(),t.$notify({title:"Success",message:"Created Successfully",type:"success",duration:2e3})})):Object(c["e"])("/fw/v2/nat",t.temp).then((function(e){t.dialogFormVisible=!1,t.getList(),t.$notify({title:"Success",message:"Created Successfully",type:"success",duration:2e3})})))}))},handleUpdate:function(t){var e=this;this.temp=Object.assign({},t),this.temp.timestamp=new Date(this.temp.timestamp),this.dialogStatus="update",this.dialogFormVisible=!0,this.$nextTick((function(){e.$refs["dataForm"].clearValidate()}))},updateData:function(){var t=this;this.$refs["dataForm"].validate((function(e){if(e){var a=Object.assign({},t.temp);a.timestamp=+new Date(a.timestamp),Object(o["e"])(a).then((function(){var e=t.list.findIndex((function(e){return e.id===t.temp.id}));t.list.splice(e,1,t.temp),t.dialogFormVisible=!1,t.$notify({title:"Success",message:"Update Successfully",type:"success",duration:2e3})}))}}))},handleDelete:function(t){var e=this,a={ip:this.$route.params.ip,forward:{protocol:t.protocol,port:t.port,toport:t.toport,toaddr:t.toaddr}};0==this.$route.params.status?Object(c["a"])("/fw/v1/nat",a).then((function(t){e.getList(),e.$notify({title:"Success",message:"Delete Successfully",type:"success",duration:2e3})})):Object(c["a"])("/fw/v2/nat",a).then((function(t){e.getList(),e.$notify({title:"Success",message:"Delete Successfully",type:"success",duration:2e3}),e.getList()}))},formatJson:function(t){return this.list.map((function(e){return t.map((function(t){return"timestamp"===t?Object(s["e"])(e[t]):e[t]}))}))},getSortClass:function(t){var e=this.listQuery.sort;return e==="+".concat(t)?"ascending":"descending"}}},f=p,m=(a("4e4f"),a("2877")),h=Object(m["a"])(f,i,n,!1,null,null,null);e["default"]=h.exports}}]);