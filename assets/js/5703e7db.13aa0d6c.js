"use strict";(self.webpackChunkmigrations=self.webpackChunkmigrations||[]).push([[396],{3905:function(e,t,n){n.d(t,{Zo:function(){return d},kt:function(){return u}});var r=n(7294);function a(e,t,n){return t in e?Object.defineProperty(e,t,{value:n,enumerable:!0,configurable:!0,writable:!0}):e[t]=n,e}function o(e,t){var n=Object.keys(e);if(Object.getOwnPropertySymbols){var r=Object.getOwnPropertySymbols(e);t&&(r=r.filter((function(t){return Object.getOwnPropertyDescriptor(e,t).enumerable}))),n.push.apply(n,r)}return n}function i(e){for(var t=1;t<arguments.length;t++){var n=null!=arguments[t]?arguments[t]:{};t%2?o(Object(n),!0).forEach((function(t){a(e,t,n[t])})):Object.getOwnPropertyDescriptors?Object.defineProperties(e,Object.getOwnPropertyDescriptors(n)):o(Object(n)).forEach((function(t){Object.defineProperty(e,t,Object.getOwnPropertyDescriptor(n,t))}))}return e}function s(e,t){if(null==e)return{};var n,r,a=function(e,t){if(null==e)return{};var n,r,a={},o=Object.keys(e);for(r=0;r<o.length;r++)n=o[r],t.indexOf(n)>=0||(a[n]=e[n]);return a}(e,t);if(Object.getOwnPropertySymbols){var o=Object.getOwnPropertySymbols(e);for(r=0;r<o.length;r++)n=o[r],t.indexOf(n)>=0||Object.prototype.propertyIsEnumerable.call(e,n)&&(a[n]=e[n])}return a}var l=r.createContext({}),c=function(e){var t=r.useContext(l),n=t;return e&&(n="function"==typeof e?e(t):i(i({},t),e)),n},d=function(e){var t=c(e.components);return r.createElement(l.Provider,{value:t},e.children)},p={inlineCode:"code",wrapper:function(e){var t=e.children;return r.createElement(r.Fragment,{},t)}},m=r.forwardRef((function(e,t){var n=e.components,a=e.mdxType,o=e.originalType,l=e.parentName,d=s(e,["components","mdxType","originalType","parentName"]),m=c(n),u=a,g=m["".concat(l,".").concat(u)]||m[u]||p[u]||o;return n?r.createElement(g,i(i({ref:t},d),{},{components:n})):r.createElement(g,i({ref:t},d))}));function u(e,t){var n=arguments,a=t&&t.mdxType;if("string"==typeof e||a){var o=n.length,i=new Array(o);i[0]=m;var s={};for(var l in t)hasOwnProperty.call(t,l)&&(s[l]=t[l]);s.originalType=e,s.mdxType="string"==typeof e?e:a,i[1]=s;for(var c=2;c<o;c++)i[c]=n[c];return r.createElement.apply(null,i)}return r.createElement.apply(null,n)}m.displayName="MDXCreateElement"},1038:function(e,t,n){n.r(t),n.d(t,{frontMatter:function(){return s},contentTitle:function(){return l},metadata:function(){return c},toc:function(){return d},default:function(){return m}});var r=n(7462),a=n(3366),o=(n(7294),n(3905)),i=["components"],s={sidebar_position:1},l="Getting Started",c={unversionedId:"sql/getting-started",id:"sql/getting-started",isDocsHomePage:!1,title:"Getting Started",description:"The migrations-sql package implements migrations based on the standard database/sql",source:"@site/docs/02-sql/01_getting-started.md",sourceDirName:"02-sql",slug:"/sql/getting-started",permalink:"/migrations/sql/getting-started",editUrl:"https://github.com/jamillosantos/migrations/edit/main/website/docs/02-sql/01_getting-started.md",tags:[],version:"current",sidebarPosition:1,frontMatter:{sidebar_position:1},sidebar:"tutorialSidebar",previous:{title:"Components",permalink:"/migrations/components"},next:{title:"Creating Migrations",permalink:"/migrations/sql/creating-a-new-migration"}},d=[{value:"What do migrations look like in my project?",id:"what-do-migrations-look-like-in-my-project",children:[],level:2},{value:"How do I trigger the migrations on my service?",id:"how-do-i-trigger-the-migrations-on-my-service",children:[],level:2}],p={toc:d};function m(e){var t=e.components,n=(0,a.Z)(e,i);return(0,o.kt)("wrapper",(0,r.Z)({},p,n,{components:t,mdxType:"MDXLayout"}),(0,o.kt)("h1",{id:"getting-started"},"Getting Started"),(0,o.kt)("p",null,"The ",(0,o.kt)("inlineCode",{parentName:"p"},"migrations-sql")," package implements migrations based on the standard ",(0,o.kt)("inlineCode",{parentName:"p"},"database/sql"),"\npackage."),(0,o.kt)("h2",{id:"what-do-migrations-look-like-in-my-project"},"What do migrations look like in my project?"),(0,o.kt)("p",null,"For the SQL databases, we can write migrations as ",(0,o.kt)("inlineCode",{parentName:"p"},".sql")," files that will be stored\nin the filesystem. So, in your project that can be a ",(0,o.kt)("inlineCode",{parentName:"p"},"migrations")," directory:"),(0,o.kt)("pre",null,(0,o.kt)("code",{parentName:"pre"},"./migrations\n\u251c\u2500\u2500 20211015015442_create_customers_table.sql\n\u251c\u2500\u2500 20211015015556_add_age_to_customers_table.sql\n\u2514\u2500\u2500 20211015045556_add_birthday_to_customers_table.sql\n")),(0,o.kt)("p",null,"By default, the ",(0,o.kt)("inlineCode",{parentName:"p"},"migrations")," package will not enable the undoing of migrations. But, if you\nenable it, you would find:"),(0,o.kt)("pre",null,(0,o.kt)("code",{parentName:"pre"},"./migrations\n\u251c\u2500\u2500 20211015015442_create_customers_table.do.sql\n\u251c\u2500\u2500 20211015015442_create_customers_table.undo.sql\n\u251c\u2500\u2500 20211015015556_add_age_to_customers_table.do.sql\n\u251c\u2500\u2500 20211015015556_add_age_to_customers_table.undo.sql\n\u2514\u2500\u2500 20211015045556_add_birthday_to_customers_table.do.sql\n\u2514\u2500\u2500 20211015045556_add_birthday_to_customers_table.undo.sql\n")),(0,o.kt)("h2",{id:"how-do-i-trigger-the-migrations-on-my-service"},"How do I trigger the migrations on my service?"),(0,o.kt)("p",null,"Create a file at ",(0,o.kt)("inlineCode",{parentName:"p"},"src/pages/my-markdown-page.md"),":"),(0,o.kt)("pre",null,(0,o.kt)("code",{parentName:"pre",className:"language-mdx",metastring:'title="src/pages/my-markdown-page.md"',title:'"src/pages/my-markdown-page.md"'},"# My Markdown page\n\nThis is a Markdown page\n")),(0,o.kt)("p",null,"A new page is now available at ",(0,o.kt)("inlineCode",{parentName:"p"},"http://localhost:3000/my-markdown-page"),"."))}m.isMDXComponent=!0}}]);