"use strict";(self.webpackChunkcosmos_sdk_docs=self.webpackChunkcosmos_sdk_docs||[]).push([[1911],{3905:(e,t,n)=>{n.d(t,{Zo:()=>d,kt:()=>u});var o=n(67294);function a(e,t,n){return t in e?Object.defineProperty(e,t,{value:n,enumerable:!0,configurable:!0,writable:!0}):e[t]=n,e}function i(e,t){var n=Object.keys(e);if(Object.getOwnPropertySymbols){var o=Object.getOwnPropertySymbols(e);t&&(o=o.filter((function(t){return Object.getOwnPropertyDescriptor(e,t).enumerable}))),n.push.apply(n,o)}return n}function r(e){for(var t=1;t<arguments.length;t++){var n=null!=arguments[t]?arguments[t]:{};t%2?i(Object(n),!0).forEach((function(t){a(e,t,n[t])})):Object.getOwnPropertyDescriptors?Object.defineProperties(e,Object.getOwnPropertyDescriptors(n)):i(Object(n)).forEach((function(t){Object.defineProperty(e,t,Object.getOwnPropertyDescriptor(n,t))}))}return e}function l(e,t){if(null==e)return{};var n,o,a=function(e,t){if(null==e)return{};var n,o,a={},i=Object.keys(e);for(o=0;o<i.length;o++)n=i[o],t.indexOf(n)>=0||(a[n]=e[n]);return a}(e,t);if(Object.getOwnPropertySymbols){var i=Object.getOwnPropertySymbols(e);for(o=0;o<i.length;o++)n=i[o],t.indexOf(n)>=0||Object.prototype.propertyIsEnumerable.call(e,n)&&(a[n]=e[n])}return a}var s=o.createContext({}),p=function(e){var t=o.useContext(s),n=t;return e&&(n="function"==typeof e?e(t):r(r({},t),e)),n},d=function(e){var t=p(e.components);return o.createElement(s.Provider,{value:t},e.children)},m={inlineCode:"code",wrapper:function(e){var t=e.children;return o.createElement(o.Fragment,{},t)}},c=o.forwardRef((function(e,t){var n=e.components,a=e.mdxType,i=e.originalType,s=e.parentName,d=l(e,["components","mdxType","originalType","parentName"]),c=p(n),u=a,h=c["".concat(s,".").concat(u)]||c[u]||m[u]||i;return n?o.createElement(h,r(r({ref:t},d),{},{components:n})):o.createElement(h,r({ref:t},d))}));function u(e,t){var n=arguments,a=t&&t.mdxType;if("string"==typeof e||a){var i=n.length,r=new Array(i);r[0]=c;var l={};for(var s in t)hasOwnProperty.call(t,s)&&(l[s]=t[s]);l.originalType=e,l.mdxType="string"==typeof e?e:a,r[1]=l;for(var p=2;p<i;p++)r[p]=n[p];return o.createElement.apply(null,r)}return o.createElement.apply(null,n)}c.displayName="MDXCreateElement"},38660:(e,t,n)=>{n.r(t),n.d(t,{assets:()=>s,contentTitle:()=>r,default:()=>m,frontMatter:()=>i,metadata:()=>l,toc:()=>p});var o=n(87462),a=(n(67294),n(3905));const i={sidebar_position:1},r="Keepers",l={unversionedId:"building-modules/keeper",id:"building-modules/keeper",title:"Keepers",description:"Keepers refer to a Cosmos SDK abstraction whose role is to manage access to the subset of the state defined by various modules. Keepers are module-specific, i.e. the subset of state defined by a module can only be accessed by a keeper defined in said module. If a module needs to access the subset of state defined by another module, a reference to the second module's internal keeper needs to be passed to the first one. This is done in app.go during the instantiation of module keepers.",source:"@site/docs/building-modules/06-keeper.md",sourceDirName:"building-modules",slug:"/building-modules/keeper",permalink:"/main/building-modules/keeper",draft:!1,tags:[],version:"current",sidebarPosition:1,frontMatter:{sidebar_position:1},sidebar:"tutorialSidebar",previous:{title:"BeginBlocker and EndBlocker",permalink:"/main/building-modules/beginblock-endblock"},next:{title:"Invariants",permalink:"/main/building-modules/invariants"}},s={},p=[{value:"Motivation",id:"motivation",level:2},{value:"Type Definition",id:"type-definition",level:2},{value:"Implementing Methods",id:"implementing-methods",level:2}],d={toc:p};function m(e){let{components:t,...n}=e;return(0,a.kt)("wrapper",(0,o.Z)({},d,n,{components:t,mdxType:"MDXLayout"}),(0,a.kt)("h1",{id:"keepers"},"Keepers"),(0,a.kt)("admonition",{title:"Synopsis",type:"note"},(0,a.kt)("p",{parentName:"admonition"},(0,a.kt)("inlineCode",{parentName:"p"},"Keeper"),"s refer to a Cosmos SDK abstraction whose role is to manage access to the subset of the state defined by various modules. ",(0,a.kt)("inlineCode",{parentName:"p"},"Keeper"),"s are module-specific, i.e. the subset of state defined by a module can only be accessed by a ",(0,a.kt)("inlineCode",{parentName:"p"},"keeper")," defined in said module. If a module needs to access the subset of state defined by another module, a reference to the second module's internal ",(0,a.kt)("inlineCode",{parentName:"p"},"keeper")," needs to be passed to the first one. This is done in ",(0,a.kt)("inlineCode",{parentName:"p"},"app.go")," during the instantiation of module keepers.")),(0,a.kt)("admonition",{title:"Pre-requisite Readings",type:"note"},(0,a.kt)("ul",{parentName:"admonition"},(0,a.kt)("li",{parentName:"ul"},(0,a.kt)("a",{parentName:"li",href:"/main/building-modules/intro"},"Introduction to Cosmos SDK Modules")))),(0,a.kt)("h2",{id:"motivation"},"Motivation"),(0,a.kt)("p",null,"The Cosmos SDK is a framework that makes it easy for developers to build complex decentralized applications from scratch, mainly by composing modules together. As the ecosystem of open-source modules for the Cosmos SDK expands, it will become increasingly likely that some of these modules contain vulnerabilities, as a result of the negligence or malice of their developer."),(0,a.kt)("p",null,"The Cosmos SDK adopts an ",(0,a.kt)("a",{parentName:"p",href:"/main/core/ocap"},"object-capabilities-based approach")," to help developers better protect their application from unwanted inter-module interactions, and ",(0,a.kt)("inlineCode",{parentName:"p"},"keeper"),"s are at the core of this approach. A ",(0,a.kt)("inlineCode",{parentName:"p"},"keeper")," can be considered quite literally to be the gatekeeper of a module's store(s). Each store (typically an ",(0,a.kt)("a",{parentName:"p",href:"/main/core/store#iavl-store"},(0,a.kt)("inlineCode",{parentName:"a"},"IAVL")," Store"),") defined within a module comes with a ",(0,a.kt)("inlineCode",{parentName:"p"},"storeKey"),", which grants unlimited access to it. The module's ",(0,a.kt)("inlineCode",{parentName:"p"},"keeper")," holds this ",(0,a.kt)("inlineCode",{parentName:"p"},"storeKey")," (which should otherwise remain unexposed), and defines ",(0,a.kt)("a",{parentName:"p",href:"#implementing-methods"},"methods")," for reading and writing to the store(s)."),(0,a.kt)("p",null,"The core idea behind the object-capabilities approach is to only reveal what is necessary to get the work done. In practice, this means that instead of handling permissions of modules through access-control lists, module ",(0,a.kt)("inlineCode",{parentName:"p"},"keeper"),"s are passed a reference to the specific instance of the other modules' ",(0,a.kt)("inlineCode",{parentName:"p"},"keeper"),"s that they need to access (this is done in the ",(0,a.kt)("a",{parentName:"p",href:"/main/basics/app-anatomy#constructor-function"},"application's constructor function"),"). As a consequence, a module can only interact with the subset of state defined in another module via the methods exposed by the instance of the other module's ",(0,a.kt)("inlineCode",{parentName:"p"},"keeper"),". This is a great way for developers to control the interactions that their own module can have with modules developed by external developers."),(0,a.kt)("h2",{id:"type-definition"},"Type Definition"),(0,a.kt)("p",null,(0,a.kt)("inlineCode",{parentName:"p"},"keeper"),"s are generally implemented in a ",(0,a.kt)("inlineCode",{parentName:"p"},"/keeper/keeper.go")," file located in the module's folder. By convention, the type ",(0,a.kt)("inlineCode",{parentName:"p"},"keeper")," of a module is simply named ",(0,a.kt)("inlineCode",{parentName:"p"},"Keeper")," and usually follows the following structure:"),(0,a.kt)("pre",null,(0,a.kt)("code",{parentName:"pre",className:"language-go"},"type Keeper struct {\n    // External keepers, if any\n\n    // Store key(s)\n\n    // codec\n\n    // authority \n}\n")),(0,a.kt)("p",null,"For example, here is the type definition of the ",(0,a.kt)("inlineCode",{parentName:"p"},"keeper")," from the ",(0,a.kt)("inlineCode",{parentName:"p"},"staking")," module:"),(0,a.kt)("pre",null,(0,a.kt)("code",{parentName:"pre",className:"language-go",metastring:"reference",reference:!0},"https://github.com/cosmos/cosmos-sdk/blob/v0.50.0-alpha.0/x/staking/keeper/keeper.go#L23-L31\n")),(0,a.kt)("p",null,"Let us go through the different parameters:"),(0,a.kt)("ul",null,(0,a.kt)("li",{parentName:"ul"},"An expected ",(0,a.kt)("inlineCode",{parentName:"li"},"keeper")," is a ",(0,a.kt)("inlineCode",{parentName:"li"},"keeper")," external to a module that is required by the internal ",(0,a.kt)("inlineCode",{parentName:"li"},"keeper")," of said module. External ",(0,a.kt)("inlineCode",{parentName:"li"},"keeper"),"s are listed in the internal ",(0,a.kt)("inlineCode",{parentName:"li"},"keeper"),"'s type definition as interfaces. These interfaces are themselves defined in an ",(0,a.kt)("inlineCode",{parentName:"li"},"expected_keepers.go")," file in the root of the module's folder. In this context, interfaces are used to reduce the number of dependencies, as well as to facilitate the maintenance of the module itself."),(0,a.kt)("li",{parentName:"ul"},(0,a.kt)("inlineCode",{parentName:"li"},"storeKey"),"s grant access to the store(s) of the ",(0,a.kt)("a",{parentName:"li",href:"/main/core/store"},"multistore")," managed by the module. They should always remain unexposed to external modules."),(0,a.kt)("li",{parentName:"ul"},(0,a.kt)("inlineCode",{parentName:"li"},"cdc")," is the ",(0,a.kt)("a",{parentName:"li",href:"/main/core/encoding"},"codec")," used to marshall and unmarshall structs to/from ",(0,a.kt)("inlineCode",{parentName:"li"},"[]byte"),". The ",(0,a.kt)("inlineCode",{parentName:"li"},"cdc")," can be any of ",(0,a.kt)("inlineCode",{parentName:"li"},"codec.BinaryCodec"),", ",(0,a.kt)("inlineCode",{parentName:"li"},"codec.JSONCodec")," or ",(0,a.kt)("inlineCode",{parentName:"li"},"codec.Codec")," based on your requirements. It can be either a proto or amino codec as long as they implement these interfaces. "),(0,a.kt)("li",{parentName:"ul"},"The authority listed is a module account or user account that has the right to change module level parameters. Previously this was handled by the param module, which has been deprecated.")),(0,a.kt)("p",null,"Of course, it is possible to define different types of internal ",(0,a.kt)("inlineCode",{parentName:"p"},"keeper"),"s for the same module (e.g. a read-only ",(0,a.kt)("inlineCode",{parentName:"p"},"keeper"),"). Each type of ",(0,a.kt)("inlineCode",{parentName:"p"},"keeper")," comes with its own constructor function, which is called from the ",(0,a.kt)("a",{parentName:"p",href:"/main/basics/app-anatomy"},"application's constructor function"),". This is where ",(0,a.kt)("inlineCode",{parentName:"p"},"keeper"),"s are instantiated, and where developers make sure to pass correct instances of modules' ",(0,a.kt)("inlineCode",{parentName:"p"},"keeper"),"s to other modules that require them."),(0,a.kt)("h2",{id:"implementing-methods"},"Implementing Methods"),(0,a.kt)("p",null,(0,a.kt)("inlineCode",{parentName:"p"},"Keeper"),"s primarily expose getter and setter methods for the store(s) managed by their module. These methods should remain as simple as possible and strictly be limited to getting or setting the requested value, as validity checks should have already been performed by the ",(0,a.kt)("a",{parentName:"p",href:"/main/building-modules/msg-services"},(0,a.kt)("inlineCode",{parentName:"a"},"Msg")," server")," when ",(0,a.kt)("inlineCode",{parentName:"p"},"keeper"),"s' methods are called."),(0,a.kt)("p",null,"Typically, a ",(0,a.kt)("em",{parentName:"p"},"getter")," method will have the following signature"),(0,a.kt)("pre",null,(0,a.kt)("code",{parentName:"pre",className:"language-go"},"func (k Keeper) Get(ctx context.Context, key string) returnType\n")),(0,a.kt)("p",null,"and the method will go through the following steps:"),(0,a.kt)("ol",null,(0,a.kt)("li",{parentName:"ol"},"Retrieve the appropriate store from the ",(0,a.kt)("inlineCode",{parentName:"li"},"ctx")," using the ",(0,a.kt)("inlineCode",{parentName:"li"},"storeKey"),". This is done through the ",(0,a.kt)("inlineCode",{parentName:"li"},"KVStore(storeKey sdk.StoreKey)")," method of the ",(0,a.kt)("inlineCode",{parentName:"li"},"ctx"),". Then it's preferred to use the ",(0,a.kt)("inlineCode",{parentName:"li"},"prefix.Store")," to access only the desired limited subset of the store for convenience and safety."),(0,a.kt)("li",{parentName:"ol"},"If it exists, get the ",(0,a.kt)("inlineCode",{parentName:"li"},"[]byte")," value stored at location ",(0,a.kt)("inlineCode",{parentName:"li"},"[]byte(key)")," using the ",(0,a.kt)("inlineCode",{parentName:"li"},"Get(key []byte)")," method of the store."),(0,a.kt)("li",{parentName:"ol"},"Unmarshall the retrieved value from ",(0,a.kt)("inlineCode",{parentName:"li"},"[]byte")," to ",(0,a.kt)("inlineCode",{parentName:"li"},"returnType")," using the codec ",(0,a.kt)("inlineCode",{parentName:"li"},"cdc"),". Return the value.")),(0,a.kt)("p",null,"Similarly, a ",(0,a.kt)("em",{parentName:"p"},"setter")," method will have the following signature"),(0,a.kt)("pre",null,(0,a.kt)("code",{parentName:"pre",className:"language-go"},"func (k Keeper) Set(ctx context.Context, key string, value valueType)\n")),(0,a.kt)("p",null,"and the method will go through the following steps:"),(0,a.kt)("ol",null,(0,a.kt)("li",{parentName:"ol"},"Retrieve the appropriate store from the ",(0,a.kt)("inlineCode",{parentName:"li"},"ctx")," using the ",(0,a.kt)("inlineCode",{parentName:"li"},"storeKey"),". This is done through the ",(0,a.kt)("inlineCode",{parentName:"li"},"KVStore(storeKey sdk.StoreKey)")," method of the ",(0,a.kt)("inlineCode",{parentName:"li"},"ctx"),". It's preferred to use the ",(0,a.kt)("inlineCode",{parentName:"li"},"prefix.Store")," to access only the desired limited subset of the store for convenience and safety."),(0,a.kt)("li",{parentName:"ol"},"Marshal ",(0,a.kt)("inlineCode",{parentName:"li"},"value")," to ",(0,a.kt)("inlineCode",{parentName:"li"},"[]byte")," using the codec ",(0,a.kt)("inlineCode",{parentName:"li"},"cdc"),"."),(0,a.kt)("li",{parentName:"ol"},"Set the encoded value in the store at location ",(0,a.kt)("inlineCode",{parentName:"li"},"key")," using the ",(0,a.kt)("inlineCode",{parentName:"li"},"Set(key []byte, value []byte)")," method of the store.")),(0,a.kt)("p",null,"For more, see an example of ",(0,a.kt)("inlineCode",{parentName:"p"},"keeper"),"'s ",(0,a.kt)("a",{parentName:"p",href:"https://github.com/cosmos/cosmos-sdk/blob/v0.50.0-alpha.0/x/staking/keeper/keeper.go"},"methods implementation from the ",(0,a.kt)("inlineCode",{parentName:"a"},"staking")," module"),"."),(0,a.kt)("p",null,"The ",(0,a.kt)("a",{parentName:"p",href:"/main/core/store#kvstore-and-commitkvstore-interfaces"},"module ",(0,a.kt)("inlineCode",{parentName:"a"},"KVStore"))," also provides an ",(0,a.kt)("inlineCode",{parentName:"p"},"Iterator()")," method which returns an ",(0,a.kt)("inlineCode",{parentName:"p"},"Iterator")," object to iterate over a domain of keys."),(0,a.kt)("p",null,"This is an example from the ",(0,a.kt)("inlineCode",{parentName:"p"},"auth")," module to iterate accounts:"),(0,a.kt)("pre",null,(0,a.kt)("code",{parentName:"pre",className:"language-go",metastring:"reference",reference:!0},"https://github.com/cosmos/cosmos-sdk/blob/v0.50.0-alpha.0/x/auth/keeper/account.go\n")))}m.isMDXComponent=!0}}]);