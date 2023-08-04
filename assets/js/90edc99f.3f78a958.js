"use strict";(self.webpackChunkmigrations=self.webpackChunkmigrations||[]).push([[593],{3905:function(e,t,n){n.d(t,{Zo:function(){return u},kt:function(){return d}});var r=n(7294);function o(e,t,n){return t in e?Object.defineProperty(e,t,{value:n,enumerable:!0,configurable:!0,writable:!0}):e[t]=n,e}function i(e,t){var n=Object.keys(e);if(Object.getOwnPropertySymbols){var r=Object.getOwnPropertySymbols(e);t&&(r=r.filter((function(t){return Object.getOwnPropertyDescriptor(e,t).enumerable}))),n.push.apply(n,r)}return n}function a(e){for(var t=1;t<arguments.length;t++){var n=null!=arguments[t]?arguments[t]:{};t%2?i(Object(n),!0).forEach((function(t){o(e,t,n[t])})):Object.getOwnPropertyDescriptors?Object.defineProperties(e,Object.getOwnPropertyDescriptors(n)):i(Object(n)).forEach((function(t){Object.defineProperty(e,t,Object.getOwnPropertyDescriptor(n,t))}))}return e}function c(e,t){if(null==e)return{};var n,r,o=function(e,t){if(null==e)return{};var n,r,o={},i=Object.keys(e);for(r=0;r<i.length;r++)n=i[r],t.indexOf(n)>=0||(o[n]=e[n]);return o}(e,t);if(Object.getOwnPropertySymbols){var i=Object.getOwnPropertySymbols(e);for(r=0;r<i.length;r++)n=i[r],t.indexOf(n)>=0||Object.prototype.propertyIsEnumerable.call(e,n)&&(o[n]=e[n])}return o}var l=r.createContext({}),s=function(e){var t=r.useContext(l),n=t;return e&&(n="function"==typeof e?e(t):a(a({},t),e)),n},u=function(e){var t=s(e.components);return r.createElement(l.Provider,{value:t},e.children)},p={inlineCode:"code",wrapper:function(e){var t=e.children;return r.createElement(r.Fragment,{},t)}},m=r.forwardRef((function(e,t){var n=e.components,o=e.mdxType,i=e.originalType,l=e.parentName,u=c(e,["components","mdxType","originalType","parentName"]),m=s(n),d=o,f=m["".concat(l,".").concat(d)]||m[d]||p[d]||i;return n?r.createElement(f,a(a({ref:t},u),{},{components:n})):r.createElement(f,a({ref:t},u))}));function d(e,t){var n=arguments,o=t&&t.mdxType;if("string"==typeof e||o){var i=n.length,a=new Array(i);a[0]=m;var c={};for(var l in t)hasOwnProperty.call(t,l)&&(c[l]=t[l]);c.originalType=e,c.mdxType="string"==typeof e?e:o,a[1]=c;for(var s=2;s<i;s++)a[s]=n[s];return r.createElement.apply(null,a)}return r.createElement.apply(null,n)}m.displayName="MDXCreateElement"},5976:function(e,t,n){n.r(t),n.d(t,{frontMatter:function(){return c},contentTitle:function(){return l},metadata:function(){return s},toc:function(){return u},default:function(){return m}});var r=n(7462),o=n(3366),i=(n(7294),n(3905)),a=["components"],c={sidebar_position:2,slug:"/components"},l="Components",s={unversionedId:"components",id:"components",isDocsHomePage:!1,title:"Components",description:"Source",source:"@site/docs/02-components.md",sourceDirName:".",slug:"/components",permalink:"/migrations/components",editUrl:"https://github.com/jamillosantos/migrations/edit/main/website/docs/02-components.md",tags:[],version:"current",sidebarPosition:2,frontMatter:{sidebar_position:2,slug:"/components"},sidebar:"tutorialSidebar",previous:{title:"Introduction",permalink:"/migrations/"},next:{title:"Getting Started",permalink:"/migrations/sql/getting-started"}},u=[{value:"Source",id:"source",children:[],level:2},{value:"Target",id:"target",children:[],level:2},{value:"Planner",id:"planner",children:[],level:2},{value:"Runner",id:"runner",children:[],level:2}],p={toc:u};function m(e){var t=e.components,n=(0,o.Z)(e,a);return(0,i.kt)("wrapper",(0,r.Z)({},p,n,{components:t,mdxType:"MDXLayout"}),(0,i.kt)("h1",{id:"components"},"Components"),(0,i.kt)("h2",{id:"source"},"Source"),(0,i.kt)("p",null,"Source is the component that will provide the migrations to be executed. It is\nresponsible for listing the migrations and loading the migration content."),(0,i.kt)("p",null,"For example, the ",(0,i.kt)("a",{parentName:"p",href:"https://github.com/jamillosantos/migration-sql"},"migration-sql")," package implements a source that can\nload migrations from the filesystem or from a ",(0,i.kt)("inlineCode",{parentName:"p"},"go:embed"),"."),(0,i.kt)("h2",{id:"target"},"Target"),(0,i.kt)("p",null,"Target is the component responsible for listing the migrations that were executed and storing the migration execution."),(0,i.kt)("h2",{id:"planner"},"Planner"),(0,i.kt)("p",null,"Planner is the component responsible for planning the migrations to be executed. It will receive the source and the\ntarget and will plan the migrations to be executed."),(0,i.kt)("h2",{id:"runner"},"Runner"),(0,i.kt)("p",null,"Runner is the component responsible for executing the migrations. It will receive the source and the target and a\nmigration plan (generated from the Planner) and executes it."))}m.isMDXComponent=!0}}]);