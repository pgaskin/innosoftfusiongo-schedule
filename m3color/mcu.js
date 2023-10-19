(()=>{var ye=Object.defineProperty;var de=Object.getOwnPropertySymbols;var De=Object.prototype.hasOwnProperty,Te=Object.prototype.propertyIsEnumerable;var ge=(e,t,r)=>t in e?ye(e,t,{enumerable:!0,configurable:!0,writable:!0,value:r}):e[t]=r,lt=(e,t)=>{for(var r in t||(t={}))De.call(t,r)&&ge(e,r,t[r]);if(de)for(var r of de(t))Te.call(t,r)&&ge(e,r,t[r]);return e};var we=(e,t)=>{for(var r in t)ye(e,r,{get:t[r],enumerable:!0})};var Ot=(e,t,r)=>new Promise((n,o)=>{var a=l=>{try{i(r.next(l))}catch(u){o(u)}},s=l=>{try{i(r.throw(l))}catch(u){o(u)}},i=l=>l.done?n(l.value):Promise.resolve(l.value).then(a,s);i((r=r.apply(e,t)).next())});var ue={};we(ue,{Blend:()=>It,Cam16:()=>N,Contrast:()=>U,CorePalette:()=>G,DislikeAnalyzer:()=>tt,DynamicColor:()=>C,DynamicScheme:()=>H,Hct:()=>F,MaterialDynamicColors:()=>c,QuantizerCelebi:()=>Ft,QuantizerMap:()=>Tt,QuantizerWsmeans:()=>Dt,QuantizerWu:()=>Bt,Scheme:()=>Pt,SchemeAndroid:()=>ne,SchemeContent:()=>oe,SchemeExpressive:()=>xt,SchemeFidelity:()=>ae,SchemeMonochrome:()=>se,SchemeNeutral:()=>ie,SchemeTonalSpot:()=>ce,SchemeVibrant:()=>Ct,Score:()=>Q,TemperatureCache:()=>mt,TonalPalette:()=>A,ViewingConditions:()=>V,alphaFromArgb:()=>Mt,applyTheme:()=>Ue,argbFromHex:()=>Ve,argbFromLab:()=>Jt,argbFromLinrgb:()=>vt,argbFromLstar:()=>Xt,argbFromRgb:()=>pt,argbFromRgba:()=>Re,argbFromXyz:()=>zt,blueFromArgb:()=>it,clampDouble:()=>ht,clampInt:()=>_t,customColor:()=>be,delinearized:()=>ot,differenceDegrees:()=>bt,greenFromArgb:()=>st,hexFromArgb:()=>Ut,isOpaque:()=>Fe,labFromArgb:()=>kt,lerp:()=>nt,linearized:()=>J,lstarFromArgb:()=>dt,lstarFromY:()=>gt,matrixMultiply:()=>ft,redFromArgb:()=>at,rgbaFromArgb:()=>Se,rotationDirection:()=>Wt,sanitizeDegreesDouble:()=>v,sanitizeDegreesInt:()=>ut,signum:()=>L,sourceColorFromImage:()=>le,themeFromImage:()=>Ne,themeFromSourceColor:()=>Ce,whitePointD65:()=>Kt,xyzFromArgb:()=>xe,yFromLstar:()=>$});function L(e){return e<0?-1:e===0?0:1}function nt(e,t,r){return(1-r)*e+r*t}function _t(e,t,r){return r<e?e:r>t?t:r}function ht(e,t,r){return r<e?e:r>t?t:r}function ut(e){return e=e%360,e<0&&(e=e+360),e}function v(e){return e=e%360,e<0&&(e=e+360),e}function Wt(e,t){return v(t-e)<=180?1:-1}function bt(e,t){return 180-Math.abs(Math.abs(e-t)-180)}function ft(e,t){let r=e[0]*t[0][0]+e[1]*t[0][1]+e[2]*t[0][2],n=e[0]*t[1][0]+e[1]*t[1][1]+e[2]*t[1][2],o=e[0]*t[2][0]+e[1]*t[2][1]+e[2]*t[2][2];return[r,n,o]}var Pe=[[.41233895,.35762064,.18051042],[.2126,.7152,.0722],[.01932141,.11916382,.95034478]],Be=[[3.2413774792388685,-1.5376652402851851,-.49885366846268053],[-.9691452513005321,1.8758853451067872,.04156585616912061],[.05562093689691305,-.20395524564742123,1.0571799111220335]],$t=[95.047,100,108.883];function pt(e,t,r){return(255<<24|(e&255)<<16|(t&255)<<8|r&255)>>>0}function vt(e){let t=ot(e[0]),r=ot(e[1]),n=ot(e[2]);return pt(t,r,n)}function Mt(e){return e>>24&255}function at(e){return e>>16&255}function st(e){return e>>8&255}function it(e){return e&255}function Fe(e){return Mt(e)>=255}function zt(e,t,r){let n=Be,o=n[0][0]*e+n[0][1]*t+n[0][2]*r,a=n[1][0]*e+n[1][1]*t+n[1][2]*r,s=n[2][0]*e+n[2][1]*t+n[2][2]*r,i=ot(o),l=ot(a),u=ot(s);return pt(i,l,u)}function xe(e){let t=J(at(e)),r=J(st(e)),n=J(it(e));return ft([t,r,n],Pe)}function Jt(e,t,r){let n=$t,o=(e+16)/116,a=t/500+o,s=o-r/200,i=Lt(a),l=Lt(o),u=Lt(s),h=i*n[0],p=l*n[1],y=u*n[2];return zt(h,p,y)}function kt(e){let t=J(at(e)),r=J(st(e)),n=J(it(e)),o=Pe,a=o[0][0]*t+o[0][1]*r+o[0][2]*n,s=o[1][0]*t+o[1][1]*r+o[1][2]*n,i=o[2][0]*t+o[2][1]*r+o[2][2]*n,l=$t,u=a/l[0],h=s/l[1],p=i/l[2],y=At(u),m=At(h),d=At(p),f=116*m-16,b=500*(y-m),k=200*(m-d);return[f,b,k]}function Xt(e){let t=$(e),r=ot(t);return pt(r,r,r)}function dt(e){let t=xe(e)[1];return 116*At(t/100)-16}function $(e){return 100*Lt((e+16)/116)}function gt(e){return At(e/100)*116-16}function J(e){let t=e/255;return t<=.040449936?t/12.92*100:Math.pow((t+.055)/1.055,2.4)*100}function ot(e){let t=e/100,r=0;return t<=.0031308?r=t*12.92:r=1.055*Math.pow(t,1/2.4)-.055,_t(0,255,Math.round(r*255))}function Kt(){return $t}function Se(e){let t=at(e),r=st(e),n=it(e),o=Mt(e);return{r:t,g:r,b:n,a:o}}function Re({r:e,g:t,b:r,a:n}){let o=Et(e),a=Et(t),s=Et(r);return Et(n)<<24|o<<16|a<<8|s}function Et(e){return e<0?0:e>255?255:e}function At(e){let t=.008856451679035631,r=24389/27;return e>t?Math.pow(e,1/3):(r*e+16)/116}function Lt(e){let t=.008856451679035631,r=24389/27,n=e*e*e;return n>t?n:(116*e-16)/r}var V=class e{static make(t=Kt(),r=200/Math.PI*$(50)/100,n=50,o=2,a=!1){let s=t,i=s[0]*.401288+s[1]*.650173+s[2]*-.051461,l=s[0]*-.250268+s[1]*1.204414+s[2]*.045854,u=s[0]*-.002079+s[1]*.048952+s[2]*.953127,h=.8+o/10,p=h>=.9?nt(.59,.69,(h-.9)*10):nt(.525,.59,(h-.8)*10),y=a?1:h*(1-1/3.6*Math.exp((-r-42)/92));y=y>1?1:y<0?0:y;let m=h,d=[y*(100/i)+1-y,y*(100/l)+1-y,y*(100/u)+1-y],f=1/(5*r+1),b=f*f*f*f,k=1-b,g=b*r+.1*k*k*Math.cbrt(5*r),P=$(n)/t[1],T=1.48+Math.sqrt(P),I=.725/Math.pow(P,.2),S=I,x=[Math.pow(g*d[0]*i/100,.42),Math.pow(g*d[1]*l/100,.42),Math.pow(g*d[2]*u/100,.42)],M=[400*x[0]/(x[0]+27.13),400*x[1]/(x[1]+27.13),400*x[2]/(x[2]+27.13)],B=(2*M[0]+M[1]+.05*M[2])*I;return new e(P,B,I,S,p,m,d,g,Math.pow(g,.25),T)}constructor(t,r,n,o,a,s,i,l,u,h){this.n=t,this.aw=r,this.nbb=n,this.ncb=o,this.c=a,this.nc=s,this.rgbD=i,this.fl=l,this.fLRoot=u,this.z=h}};V.DEFAULT=V.make();var N=class e{constructor(t,r,n,o,a,s,i,l,u){this.hue=t,this.chroma=r,this.j=n,this.q=o,this.m=a,this.s=s,this.jstar=i,this.astar=l,this.bstar=u}distance(t){let r=this.jstar-t.jstar,n=this.astar-t.astar,o=this.bstar-t.bstar,a=Math.sqrt(r*r+n*n+o*o);return 1.41*Math.pow(a,.63)}static fromInt(t){return e.fromIntInViewingConditions(t,V.DEFAULT)}static fromIntInViewingConditions(t,r){let n=(t&16711680)>>16,o=(t&65280)>>8,a=t&255,s=J(n),i=J(o),l=J(a),u=.41233895*s+.35762064*i+.18051042*l,h=.2126*s+.7152*i+.0722*l,p=.01932141*s+.11916382*i+.95034478*l,y=.401288*u+.650173*h-.051461*p,m=-.250268*u+1.204414*h+.045854*p,d=-.002079*u+.048952*h+.953127*p,f=r.rgbD[0]*y,b=r.rgbD[1]*m,k=r.rgbD[2]*d,g=Math.pow(r.fl*Math.abs(f)/100,.42),P=Math.pow(r.fl*Math.abs(b)/100,.42),T=Math.pow(r.fl*Math.abs(k)/100,.42),I=L(f)*400*g/(g+27.13),S=L(b)*400*P/(P+27.13),x=L(k)*400*T/(T+27.13),M=(11*I+-12*S+x)/11,B=(I+S-2*x)/9,w=(20*I+20*S+21*x)/20,q=(40*I+20*S+x)/20,Y=Math.atan2(B,M)*180/Math.PI,E=Y<0?Y+360:Y>=360?Y-360:Y,ct=E*Math.PI/180,St=q*r.nbb,rt=100*Math.pow(St/r.aw,r.c*r.z),Rt=4/r.c*Math.sqrt(rt/100)*(r.aw+4)*r.fLRoot,Gt=E<20.14?E+360:E,qt=.25*(Math.cos(Gt*Math.PI/180+2)+3.8),Yt=5e4/13*qt*r.nc*r.ncb*Math.sqrt(M*M+B*B)/(w+.305),Ht=Math.pow(Yt,.9)*Math.pow(1.64-Math.pow(.29,r.n),.73),me=Ht*Math.sqrt(rt/100),fe=me*r.fLRoot,Ae=50*Math.sqrt(Ht*r.c/(r.aw+4)),Me=(1+100*.007)*rt/(1+.007*rt),pe=1/.0228*Math.log(1+.0228*fe),ke=pe*Math.cos(ct),Ie=pe*Math.sin(ct);return new e(E,me,rt,Rt,fe,Ae,Me,ke,Ie)}static fromJch(t,r,n){return e.fromJchInViewingConditions(t,r,n,V.DEFAULT)}static fromJchInViewingConditions(t,r,n,o){let a=4/o.c*Math.sqrt(t/100)*(o.aw+4)*o.fLRoot,s=r*o.fLRoot,i=r/Math.sqrt(t/100),l=50*Math.sqrt(i*o.c/(o.aw+4)),u=n*Math.PI/180,h=(1+100*.007)*t/(1+.007*t),p=1/.0228*Math.log(1+.0228*s),y=p*Math.cos(u),m=p*Math.sin(u);return new e(n,r,t,a,s,l,h,y,m)}static fromUcs(t,r,n){return e.fromUcsInViewingConditions(t,r,n,V.DEFAULT)}static fromUcsInViewingConditions(t,r,n,o){let a=r,s=n,i=Math.sqrt(a*a+s*s),u=(Math.exp(i*.0228)-1)/.0228/o.fLRoot,h=Math.atan2(s,a)*(180/Math.PI);h<0&&(h+=360);let p=t/(1-(t-100)*.007);return e.fromJchInViewingConditions(p,u,h,o)}toInt(){return this.viewed(V.DEFAULT)}viewed(t){let r=this.chroma===0||this.j===0?0:this.chroma/Math.sqrt(this.j/100),n=Math.pow(r/Math.pow(1.64-Math.pow(.29,t.n),.73),1/.9),o=this.hue*Math.PI/180,a=.25*(Math.cos(o+2)+3.8),s=t.aw*Math.pow(this.j/100,1/t.c/t.z),i=a*(5e4/13)*t.nc*t.ncb,l=s/t.nbb,u=Math.sin(o),h=Math.cos(o),p=23*(l+.305)*n/(23*i+11*n*h+108*n*u),y=p*h,m=p*u,d=(460*l+451*y+288*m)/1403,f=(460*l-891*y-261*m)/1403,b=(460*l-220*y-6300*m)/1403,k=Math.max(0,27.13*Math.abs(d)/(400-Math.abs(d))),g=L(d)*(100/t.fl)*Math.pow(k,1/.42),P=Math.max(0,27.13*Math.abs(f)/(400-Math.abs(f))),T=L(f)*(100/t.fl)*Math.pow(P,1/.42),I=Math.max(0,27.13*Math.abs(b)/(400-Math.abs(b))),S=L(b)*(100/t.fl)*Math.pow(I,1/.42),x=g/t.rgbD[0],M=T/t.rgbD[1],B=S/t.rgbD[2],w=1.86206786*x-1.01125463*M+.14918677*B,q=.38752654*x+.62144744*M-.00897398*B,W=-.0158415*x-.03412294*M+1.04996444*B;return zt(w,q,W)}static fromXyzInViewingConditions(t,r,n,o){let a=.401288*t+.650173*r-.051461*n,s=-.250268*t+1.204414*r+.045854*n,i=-.002079*t+.048952*r+.953127*n,l=o.rgbD[0]*a,u=o.rgbD[1]*s,h=o.rgbD[2]*i,p=Math.pow(o.fl*Math.abs(l)/100,.42),y=Math.pow(o.fl*Math.abs(u)/100,.42),m=Math.pow(o.fl*Math.abs(h)/100,.42),d=L(l)*400*p/(p+27.13),f=L(u)*400*y/(y+27.13),b=L(h)*400*m/(m+27.13),k=(11*d+-12*f+b)/11,g=(d+f-2*b)/9,P=(20*d+20*f+21*b)/20,T=(40*d+20*f+b)/20,S=Math.atan2(g,k)*180/Math.PI,x=S<0?S+360:S>=360?S-360:S,M=x*Math.PI/180,B=T*o.nbb,w=100*Math.pow(B/o.aw,o.c*o.z),q=4/o.c*Math.sqrt(w/100)*(o.aw+4)*o.fLRoot,W=x<20.14?x+360:x,Y=1/4*(Math.cos(W*Math.PI/180+2)+3.8),ct=5e4/13*Y*o.nc*o.ncb*Math.sqrt(k*k+g*g)/(P+.305),St=Math.pow(ct,.9)*Math.pow(1.64-Math.pow(.29,o.n),.73),rt=St*Math.sqrt(w/100),Rt=rt*o.fLRoot,Gt=50*Math.sqrt(St*o.c/(o.aw+4)),qt=(1+100*.007)*w/(1+.007*w),jt=Math.log(1+.0228*Rt)/.0228,Yt=jt*Math.cos(M),Ht=jt*Math.sin(M);return new e(x,rt,w,q,Rt,Gt,qt,Yt,Ht)}xyzInViewingConditions(t){let r=this.chroma===0||this.j===0?0:this.chroma/Math.sqrt(this.j/100),n=Math.pow(r/Math.pow(1.64-Math.pow(.29,t.n),.73),1/.9),o=this.hue*Math.PI/180,a=.25*(Math.cos(o+2)+3.8),s=t.aw*Math.pow(this.j/100,1/t.c/t.z),i=a*(5e4/13)*t.nc*t.ncb,l=s/t.nbb,u=Math.sin(o),h=Math.cos(o),p=23*(l+.305)*n/(23*i+11*n*h+108*n*u),y=p*h,m=p*u,d=(460*l+451*y+288*m)/1403,f=(460*l-891*y-261*m)/1403,b=(460*l-220*y-6300*m)/1403,k=Math.max(0,27.13*Math.abs(d)/(400-Math.abs(d))),g=L(d)*(100/t.fl)*Math.pow(k,1/.42),P=Math.max(0,27.13*Math.abs(f)/(400-Math.abs(f))),T=L(f)*(100/t.fl)*Math.pow(P,1/.42),I=Math.max(0,27.13*Math.abs(b)/(400-Math.abs(b))),S=L(b)*(100/t.fl)*Math.pow(I,1/.42),x=g/t.rgbD[0],M=T/t.rgbD[1],B=S/t.rgbD[2],w=1.86206786*x-1.01125463*M+.14918677*B,q=.38752654*x+.62144744*M-.00897398*B,W=-.0158415*x-.03412294*M+1.04996444*B;return[w,q,W]}};var K=class e{static sanitizeRadians(t){return(t+Math.PI*8)%(Math.PI*2)}static trueDelinearized(t){let r=t/100,n=0;return r<=.0031308?n=r*12.92:n=1.055*Math.pow(r,1/2.4)-.055,n*255}static chromaticAdaptation(t){let r=Math.pow(Math.abs(t),.42);return L(t)*400*r/(r+27.13)}static hueOf(t){let r=ft(t,e.SCALED_DISCOUNT_FROM_LINRGB),n=e.chromaticAdaptation(r[0]),o=e.chromaticAdaptation(r[1]),a=e.chromaticAdaptation(r[2]),s=(11*n+-12*o+a)/11,i=(n+o-2*a)/9;return Math.atan2(i,s)}static areInCyclicOrder(t,r,n){let o=e.sanitizeRadians(r-t),a=e.sanitizeRadians(n-t);return o<a}static intercept(t,r,n){return(r-t)/(n-t)}static lerpPoint(t,r,n){return[t[0]+(n[0]-t[0])*r,t[1]+(n[1]-t[1])*r,t[2]+(n[2]-t[2])*r]}static setCoordinate(t,r,n,o){let a=e.intercept(t[o],r,n[o]);return e.lerpPoint(t,a,n)}static isBounded(t){return 0<=t&&t<=100}static nthVertex(t,r){let n=e.Y_FROM_LINRGB[0],o=e.Y_FROM_LINRGB[1],a=e.Y_FROM_LINRGB[2],s=r%4<=1?0:100,i=r%2===0?0:100;if(r<4){let l=s,u=i,h=(t-l*o-u*a)/n;return e.isBounded(h)?[h,l,u]:[-1,-1,-1]}else if(r<8){let l=s,u=i,h=(t-u*n-l*a)/o;return e.isBounded(h)?[u,h,l]:[-1,-1,-1]}else{let l=s,u=i,h=(t-l*n-u*o)/a;return e.isBounded(h)?[l,u,h]:[-1,-1,-1]}}static bisectToSegment(t,r){let n=[-1,-1,-1],o=n,a=0,s=0,i=!1,l=!0;for(let u=0;u<12;u++){let h=e.nthVertex(t,u);if(h[0]<0)continue;let p=e.hueOf(h);if(!i){n=h,o=h,a=p,s=p,i=!0;continue}(l||e.areInCyclicOrder(a,p,s))&&(l=!1,e.areInCyclicOrder(a,r,p)?(o=h,s=p):(n=h,a=p))}return[n,o]}static midpoint(t,r){return[(t[0]+r[0])/2,(t[1]+r[1])/2,(t[2]+r[2])/2]}static criticalPlaneBelow(t){return Math.floor(t-.5)}static criticalPlaneAbove(t){return Math.ceil(t-.5)}static bisectToLimit(t,r){let n=e.bisectToSegment(t,r),o=n[0],a=e.hueOf(o),s=n[1];for(let i=0;i<3;i++)if(o[i]!==s[i]){let l=-1,u=255;o[i]<s[i]?(l=e.criticalPlaneBelow(e.trueDelinearized(o[i])),u=e.criticalPlaneAbove(e.trueDelinearized(s[i]))):(l=e.criticalPlaneAbove(e.trueDelinearized(o[i])),u=e.criticalPlaneBelow(e.trueDelinearized(s[i])));for(let h=0;h<8&&!(Math.abs(u-l)<=1);h++){let p=Math.floor((l+u)/2),y=e.CRITICAL_PLANES[p],m=e.setCoordinate(o,y,s,i),d=e.hueOf(m);e.areInCyclicOrder(a,r,d)?(s=m,u=p):(o=m,a=d,l=p)}}return e.midpoint(o,s)}static inverseChromaticAdaptation(t){let r=Math.abs(t),n=Math.max(0,27.13*r/(400-r));return L(t)*Math.pow(n,1/.42)}static findResultByJ(t,r,n){let o=Math.sqrt(n)*11,a=V.DEFAULT,s=1/Math.pow(1.64-Math.pow(.29,a.n),.73),l=.25*(Math.cos(t+2)+3.8)*(5e4/13)*a.nc*a.ncb,u=Math.sin(t),h=Math.cos(t);for(let p=0;p<5;p++){let y=o/100,m=r===0||o===0?0:r/Math.sqrt(y),d=Math.pow(m*s,1/.9),b=a.aw*Math.pow(y,1/a.c/a.z)/a.nbb,k=23*(b+.305)*d/(23*l+11*d*h+108*d*u),g=k*h,P=k*u,T=(460*b+451*g+288*P)/1403,I=(460*b-891*g-261*P)/1403,S=(460*b-220*g-6300*P)/1403,x=e.inverseChromaticAdaptation(T),M=e.inverseChromaticAdaptation(I),B=e.inverseChromaticAdaptation(S),w=ft([x,M,B],e.LINRGB_FROM_SCALED_DISCOUNT);if(w[0]<0||w[1]<0||w[2]<0)return 0;let q=e.Y_FROM_LINRGB[0],W=e.Y_FROM_LINRGB[1],Y=e.Y_FROM_LINRGB[2],E=q*w[0]+W*w[1]+Y*w[2];if(E<=0)return 0;if(p===4||Math.abs(E-n)<.002)return w[0]>100.01||w[1]>100.01||w[2]>100.01?0:vt(w);o=o-(E-n)*o/(2*E)}return 0}static solveToInt(t,r,n){if(r<1e-4||n<1e-4||n>99.9999)return Xt(n);t=v(t);let o=t/180*Math.PI,a=$(n),s=e.findResultByJ(o,r,a);if(s!==0)return s;let i=e.bisectToLimit(a,o);return vt(i)}static solveToCam(t,r,n){return N.fromInt(e.solveToInt(t,r,n))}};K.SCALED_DISCOUNT_FROM_LINRGB=[[.001200833568784504,.002389694492170889,.0002795742885861124],[.0005891086651375999,.0029785502573438758,.0003270666104008398],[.00010146692491640572,.0005364214359186694,.0032979401770712076]];K.LINRGB_FROM_SCALED_DISCOUNT=[[1373.2198709594231,-1100.4251190754821,-7.278681089101213],[-271.815969077903,559.6580465940733,-32.46047482791194],[1.9622899599665666,-57.173814538844006,308.7233197812385]];K.Y_FROM_LINRGB=[.2126,.7152,.0722];K.CRITICAL_PLANES=[.015176349177441876,.045529047532325624,.07588174588720938,.10623444424209313,.13658714259697685,.16693984095186062,.19729253930674434,.2276452376616281,.2579979360165119,.28835063437139563,.3188300904430532,.350925934958123,.3848314933096426,.42057480301049466,.458183274052838,.4976837250274023,.5391024159806381,.5824650784040898,.6277969426914107,.6751227633498623,.7244668422128921,.775853049866786,.829304845476233,.8848452951698498,.942497089126609,1.0022825574869039,1.0642236851973577,1.1283421258858297,1.1946592148522128,1.2631959812511864,1.3339731595349034,1.407011200216447,1.4823302800086415,1.5599503113873272,1.6398909516233677,1.7221716113234105,1.8068114625156377,1.8938294463134073,1.9832442801866852,2.075074464868551,2.1693382909216234,2.2660538449872063,2.36523901573795,2.4669114995532007,2.5710888059345764,2.6777882626779785,2.7870270208169257,2.898822059350997,3.0131901897720907,3.1301480604002863,3.2497121605402226,3.3718988244681087,3.4967242352587946,3.624204428461639,3.754355295633311,3.887192587735158,4.022731918402185,4.160988767090289,4.301978482107941,4.445716283538092,4.592217266055746,4.741496401646282,4.893568542229298,5.048448422192488,5.20615066083972,5.3666897647573375,5.5300801301023865,5.696336044816294,5.865471690767354,6.037501145825082,6.212438385869475,6.390297286737924,6.571091626112461,6.7548350853498045,6.941541251256611,7.131223617812143,7.323895587840543,7.5195704746346665,7.7182615035334345,7.919981813454504,8.124744458384042,8.332562408825165,8.543448553206703,8.757415699253682,8.974476575321063,9.194643831691977,9.417930041841839,9.644347703669503,9.873909240696694,10.106627003236781,10.342513269534024,10.58158024687427,10.8238400726681,11.069304815507364,11.317986476196008,11.569896988756009,11.825048221409341,12.083451977536606,12.345119996613247,12.610063955123938,12.878295467455942,13.149826086772048,13.42466730586372,13.702830557985108,13.984327217668513,14.269168601521828,14.55736596900856,14.848930523210871,15.143873411576273,15.44220572664832,15.743938506781891,16.04908273684337,16.35764934889634,16.66964922287304,16.985093187232053,17.30399201960269,17.62635644741625,17.95219714852476,18.281524751807332,18.614349837764564,18.95068293910138,19.290534541298456,19.633915083172692,19.98083495742689,20.331304511189067,20.685334046541502,21.042933821039977,21.404114048223256,21.76888489811322,22.137256497705877,22.50923893145328,22.884842241736916,23.264076429332462,23.6469514538663,24.033477234264016,24.42366364919083,24.817520537484558,25.21505769858089,25.61628489293138,26.021211842414342,26.429848230738664,26.842203703840827,27.258287870275353,27.678110301598522,28.10168053274597,28.529008062403893,28.96010235337422,29.39497283293396,29.83362889318845,30.276079891419332,30.722335150426627,31.172403958865512,31.62629557157785,32.08401920991837,32.54558406207592,33.010999283389665,33.4802739966603,33.953417292456834,34.430438229418264,34.911345834551085,35.39614910352207,35.88485700094671,36.37747846067349,36.87402238606382,37.37449765026789,37.87891309649659,38.38727753828926,38.89959975977785,39.41588851594697,39.93615253289054,40.460400508064545,40.98864111053629,41.520882981230194,42.05713473317016,42.597404951718396,43.141702194811224,43.6900349931913,44.24241185063697,44.798841244188324,45.35933162437017,45.92389141541209,46.49252901546552,47.065252796817916,47.64207110610409,48.22299226451468,48.808024568002054,49.3971762874833,49.9904556690408,50.587870934119984,51.189430279724725,51.79514187861014,52.40501387947288,53.0190544071392,53.637271562750364,54.259673423945976,54.88626804504493,55.517063457223934,56.15206766869424,56.79128866487574,57.43473440856916,58.08241284012621,58.734331877617365,59.39049941699807,60.05092333227251,60.715611475655585,61.38457167773311,62.057811747619894,62.7353394731159,63.417162620860914,64.10328893648692,64.79372614476921,65.48848194977529,66.18756403501224,66.89098006357258,67.59873767827808,68.31084450182222,69.02730813691093,69.74813616640164,70.47333615344107,71.20291564160104,71.93688215501312,72.67524319850172,73.41800625771542,74.16517879925733,74.9167682708136,75.67278210128072,76.43322770089146,77.1981124613393,77.96744375590167,78.74122893956174,79.51947534912904,80.30219030335869,81.08938110306934,81.88105503125999,82.67721935322541,83.4778813166706,84.28304815182372,85.09272707154808,85.90692527145302,86.72564993000343,87.54890820862819,88.3767072518277,89.2090541872801,90.04595612594655,90.88742016217518,91.73345337380438,92.58406282226491,93.43925555268066,94.29903859396902,95.16341895893969,96.03240364439274,96.9059996312159,97.78421388448044,98.6670533535366,99.55452497210776];var F=class e{static from(t,r,n){return new e(K.solveToInt(t,r,n))}static fromInt(t){return new e(t)}toInt(){return this.argb}get hue(){return this.internalHue}set hue(t){this.setInternalState(K.solveToInt(t,this.internalChroma,this.internalTone))}get chroma(){return this.internalChroma}set chroma(t){this.setInternalState(K.solveToInt(this.internalHue,t,this.internalTone))}get tone(){return this.internalTone}set tone(t){this.setInternalState(K.solveToInt(this.internalHue,this.internalChroma,t))}constructor(t){this.argb=t;let r=N.fromInt(t);this.internalHue=r.hue,this.internalChroma=r.chroma,this.internalTone=dt(t),this.argb=t}setInternalState(t){let r=N.fromInt(t);this.internalHue=r.hue,this.internalChroma=r.chroma,this.internalTone=dt(t),this.argb=t}inViewingConditions(t){let n=N.fromInt(this.toInt()).xyzInViewingConditions(t),o=N.fromXyzInViewingConditions(n[0],n[1],n[2],V.make());return e.from(o.hue,o.chroma,gt(n[1]))}};var It=class e{static harmonize(t,r){let n=F.fromInt(t),o=F.fromInt(r),a=bt(n.hue,o.hue),s=Math.min(a*.5,15),i=v(n.hue+s*Wt(n.hue,o.hue));return F.from(i,n.chroma,n.tone).toInt()}static hctHue(t,r,n){let o=e.cam16Ucs(t,r,n),a=N.fromInt(o),s=N.fromInt(t);return F.from(a.hue,s.chroma,dt(t)).toInt()}static cam16Ucs(t,r,n){let o=N.fromInt(t),a=N.fromInt(r),s=o.jstar,i=o.astar,l=o.bstar,u=a.jstar,h=a.astar,p=a.bstar,y=s+(u-s)*n,m=i+(h-i)*n,d=l+(p-l)*n;return N.fromUcs(y,m,d).toInt()}};var U=class e{static ratioOfTones(t,r){return t=ht(0,100,t),r=ht(0,100,r),e.ratioOfYs($(t),$(r))}static ratioOfYs(t,r){let n=t>r?t:r,o=n===r?t:r;return(n+5)/(o+5)}static lighter(t,r){if(t<0||t>100)return-1;let n=$(t),o=r*(n+5)-5,a=e.ratioOfYs(o,n),s=Math.abs(a-r);if(a<r&&s>.04)return-1;let i=gt(o)+.4;return i<0||i>100?-1:i}static darker(t,r){if(t<0||t>100)return-1;let n=$(t),o=(n+5)/r-5,a=e.ratioOfYs(n,o),s=Math.abs(a-r);if(a<r&&s>.04)return-1;let i=gt(o)-.4;return i<0||i>100?-1:i}static lighterUnsafe(t,r){let n=e.lighter(t,r);return n<0?100:n}static darkerUnsafe(t,r){let n=e.darker(t,r);return n<0?0:n}};var tt=class e{static isDisliked(t){let r=Math.round(t.hue)>=90&&Math.round(t.hue)<=111,n=Math.round(t.chroma)>16,o=Math.round(t.tone)<65;return r&&n&&o}static fixIfDisliked(t){return e.isDisliked(t)?F.from(t.hue,t.chroma,70):t}};var C=class e{static fromPalette(t){var r,n;return new e((r=t.name)!=null?r:"",t.palette,t.tone,(n=t.isBackground)!=null?n:!1,t.background,t.secondBackground,t.contrastCurve,t.toneDeltaPair)}constructor(t,r,n,o,a,s,i,l){if(this.name=t,this.palette=r,this.tone=n,this.isBackground=o,this.background=a,this.secondBackground=s,this.contrastCurve=i,this.toneDeltaPair=l,this.hctCache=new Map,!a&&s)throw new Error(`Color ${t} has secondBackgrounddefined, but background is not defined.`);if(!a&&i)throw new Error(`Color ${t} has contrastCurvedefined, but background is not defined.`);if(a&&!i)throw new Error(`Color ${t} has backgrounddefined, but contrastCurve is not defined.`)}getArgb(t){return this.getHct(t).toInt()}getHct(t){let r=this.hctCache.get(t);if(r!=null)return r;let n=this.getTone(t),o=this.palette(t).getHct(n);return this.hctCache.size>4&&this.hctCache.clear(),this.hctCache.set(t,o),o}getTone(t){let r=t.contrastLevel<0;if(this.toneDeltaPair){let n=this.toneDeltaPair(t),o=n.roleA,a=n.roleB,s=n.delta,i=n.polarity,l=n.stayTogether,h=this.background(t).getTone(t),p=i==="nearer"||i==="lighter"&&!t.isDark||i==="darker"&&t.isDark,y=p?o:a,m=p?a:o,d=this.name===y.name,f=t.isDark?1:-1,b=y.contrastCurve.getContrast(t.contrastLevel),k=m.contrastCurve.getContrast(t.contrastLevel),g=y.tone(t),P=U.ratioOfTones(h,g)>=b?g:e.foregroundTone(h,b),T=m.tone(t),I=U.ratioOfTones(h,T)>=k?T:e.foregroundTone(h,k);return r&&(P=e.foregroundTone(h,b),I=e.foregroundTone(h,k)),(I-P)*f>=s||(I=ht(0,100,P+s*f),(I-P)*f>=s||(P=ht(0,100,I-s*f))),50<=P&&P<60?f>0?(P=60,I=Math.max(I,P+s*f)):(P=49,I=Math.min(I,P+s*f)):50<=I&&I<60&&(l?f>0?(P=60,I=Math.max(I,P+s*f)):(P=49,I=Math.min(I,P+s*f)):f>0?I=60:I=49),d?P:I}else{let n=this.tone(t);if(this.background==null)return n;let o=this.background(t).getTone(t),a=this.contrastCurve.getContrast(t.contrastLevel);if(U.ratioOfTones(o,n)>=a||(n=e.foregroundTone(o,a)),r&&(n=e.foregroundTone(o,a)),this.isBackground&&50<=n&&n<60&&(U.ratioOfTones(49,o)>=a?n=49:n=60),this.secondBackground){let[s,i]=[this.background,this.secondBackground],[l,u]=[s(t).getTone(t),i(t).getTone(t)],[h,p]=[Math.max(l,u),Math.min(l,u)];if(U.ratioOfTones(h,n)>=a&&U.ratioOfTones(p,n)>=a)return n;let y=U.lighter(h,a),m=U.darker(p,a),d=[];return y!==-1&&d.push(y),m!==-1&&d.push(m),e.tonePrefersLightForeground(l)||e.tonePrefersLightForeground(u)?y<0?100:y:d.length===1?d[0]:m<0?0:m}return n}}static foregroundTone(t,r){let n=U.lighterUnsafe(t,r),o=U.darkerUnsafe(t,r),a=U.ratioOfTones(n,t),s=U.ratioOfTones(o,t);if(e.tonePrefersLightForeground(t)){let l=Math.abs(a-s)<.1&&a<r&&s<r;return a>=r||a>=s||l?n:o}else return s>=r||s>=a?o:n}static tonePrefersLightForeground(t){return Math.round(t)<60}static toneAllowsLightForeground(t){return Math.round(t)<=49}static enableLightForeground(t){return e.tonePrefersLightForeground(t)&&!e.toneAllowsLightForeground(t)?49:t}};var O;(function(e){e[e.MONOCHROME=0]="MONOCHROME",e[e.NEUTRAL=1]="NEUTRAL",e[e.TONAL_SPOT=2]="TONAL_SPOT",e[e.VIBRANT=3]="VIBRANT",e[e.EXPRESSIVE=4]="EXPRESSIVE",e[e.FIDELITY=5]="FIDELITY",e[e.CONTENT=6]="CONTENT",e[e.RAINBOW=7]="RAINBOW",e[e.FRUIT_SALAD=8]="FRUIT_SALAD"})(O||(O={}));var D=class{constructor(t,r,n,o){this.low=t,this.normal=r,this.medium=n,this.high=o}getContrast(t){return t<=-1?this.low:t<0?nt(this.low,this.normal,(t- -1)/1):t<.5?nt(this.normal,this.medium,(t-0)/.5):t<1?nt(this.medium,this.high,(t-.5)/.5):this.high}};var z=class{constructor(t,r,n,o,a){this.roleA=t,this.roleB=r,this.delta=n,this.polarity=o,this.stayTogether=a}};function yt(e){return e.variant===O.FIDELITY||e.variant===O.CONTENT}function R(e){return e.variant===O.MONOCHROME}function He(e,t,r,n){let o=r,a=F.from(e,t,r);if(a.chroma<t){let s=a.chroma;for(;a.chroma<t;){o+=n?-1:1;let i=F.from(e,t,o);if(s>i.chroma||Math.abs(i.chroma-t)<.4)break;let l=Math.abs(i.chroma-t),u=Math.abs(a.chroma-t);l<u&&(a=i),s=Math.max(s,i.chroma)}}return o}function Oe(e){return V.make(void 0,void 0,e.isDark?30:80,void 0,void 0)}function Zt(e,t){let r=e.inViewingConditions(Oe(t));return C.tonePrefersLightForeground(e.tone)&&!C.toneAllowsLightForeground(r.tone)?C.enableLightForeground(e.tone):C.enableLightForeground(r.tone)}var c=class e{static highestSurface(t){return t.isDark?e.surfaceBright:e.surfaceDim}};c.contentAccentToneDelta=15;c.primaryPaletteKeyColor=C.fromPalette({name:"primary_palette_key_color",palette:e=>e.primaryPalette,tone:e=>e.primaryPalette.keyColor.tone});c.secondaryPaletteKeyColor=C.fromPalette({name:"secondary_palette_key_color",palette:e=>e.secondaryPalette,tone:e=>e.secondaryPalette.keyColor.tone});c.tertiaryPaletteKeyColor=C.fromPalette({name:"tertiary_palette_key_color",palette:e=>e.tertiaryPalette,tone:e=>e.tertiaryPalette.keyColor.tone});c.neutralPaletteKeyColor=C.fromPalette({name:"neutral_palette_key_color",palette:e=>e.neutralPalette,tone:e=>e.neutralPalette.keyColor.tone});c.neutralVariantPaletteKeyColor=C.fromPalette({name:"neutral_variant_palette_key_color",palette:e=>e.neutralVariantPalette,tone:e=>e.neutralVariantPalette.keyColor.tone});c.background=C.fromPalette({name:"background",palette:e=>e.neutralPalette,tone:e=>e.isDark?6:98,isBackground:!0});c.onBackground=C.fromPalette({name:"on_background",palette:e=>e.neutralPalette,tone:e=>e.isDark?90:10,background:e=>c.background,contrastCurve:new D(3,3,4.5,7)});c.surface=C.fromPalette({name:"surface",palette:e=>e.neutralPalette,tone:e=>e.isDark?6:98,isBackground:!0});c.surfaceDim=C.fromPalette({name:"surface_dim",palette:e=>e.neutralPalette,tone:e=>e.isDark?6:87,isBackground:!0});c.surfaceBright=C.fromPalette({name:"surface_bright",palette:e=>e.neutralPalette,tone:e=>e.isDark?24:98,isBackground:!0});c.surfaceContainerLowest=C.fromPalette({name:"surface_container_lowest",palette:e=>e.neutralPalette,tone:e=>e.isDark?4:100,isBackground:!0});c.surfaceContainerLow=C.fromPalette({name:"surface_container_low",palette:e=>e.neutralPalette,tone:e=>e.isDark?10:96,isBackground:!0});c.surfaceContainer=C.fromPalette({name:"surface_container",palette:e=>e.neutralPalette,tone:e=>e.isDark?12:94,isBackground:!0});c.surfaceContainerHigh=C.fromPalette({name:"surface_container_high",palette:e=>e.neutralPalette,tone:e=>e.isDark?17:92,isBackground:!0});c.surfaceContainerHighest=C.fromPalette({name:"surface_container_highest",palette:e=>e.neutralPalette,tone:e=>e.isDark?22:90,isBackground:!0});c.onSurface=C.fromPalette({name:"on_surface",palette:e=>e.neutralPalette,tone:e=>e.isDark?90:10,background:e=>c.highestSurface(e),contrastCurve:new D(4.5,7,11,21)});c.surfaceVariant=C.fromPalette({name:"surface_variant",palette:e=>e.neutralVariantPalette,tone:e=>e.isDark?30:90,isBackground:!0});c.onSurfaceVariant=C.fromPalette({name:"on_surface_variant",palette:e=>e.neutralVariantPalette,tone:e=>e.isDark?80:30,background:e=>c.highestSurface(e),contrastCurve:new D(3,4.5,7,11)});c.inverseSurface=C.fromPalette({name:"inverse_surface",palette:e=>e.neutralPalette,tone:e=>e.isDark?90:20});c.inverseOnSurface=C.fromPalette({name:"inverse_on_surface",palette:e=>e.neutralPalette,tone:e=>e.isDark?20:95,background:e=>c.inverseSurface,contrastCurve:new D(4.5,7,11,21)});c.outline=C.fromPalette({name:"outline",palette:e=>e.neutralVariantPalette,tone:e=>e.isDark?60:50,background:e=>c.highestSurface(e),contrastCurve:new D(1.5,3,4.5,7)});c.outlineVariant=C.fromPalette({name:"outline_variant",palette:e=>e.neutralVariantPalette,tone:e=>e.isDark?30:80,background:e=>c.highestSurface(e),contrastCurve:new D(1,1,3,7)});c.shadow=C.fromPalette({name:"shadow",palette:e=>e.neutralPalette,tone:e=>0});c.scrim=C.fromPalette({name:"scrim",palette:e=>e.neutralPalette,tone:e=>0});c.surfaceTint=C.fromPalette({name:"surface_tint",palette:e=>e.primaryPalette,tone:e=>e.isDark?80:40,isBackground:!0});c.primary=C.fromPalette({name:"primary",palette:e=>e.primaryPalette,tone:e=>R(e)?e.isDark?100:0:e.isDark?80:40,isBackground:!0,background:e=>c.highestSurface(e),contrastCurve:new D(3,4.5,7,11),toneDeltaPair:e=>new z(c.primaryContainer,c.primary,15,"nearer",!1)});c.onPrimary=C.fromPalette({name:"on_primary",palette:e=>e.primaryPalette,tone:e=>R(e)?e.isDark?10:90:e.isDark?20:100,background:e=>c.primary,contrastCurve:new D(4.5,7,11,21)});c.primaryContainer=C.fromPalette({name:"primary_container",palette:e=>e.primaryPalette,tone:e=>yt(e)?Zt(e.sourceColorHct,e):R(e)?e.isDark?85:25:e.isDark?30:90,isBackground:!0,background:e=>c.highestSurface(e),contrastCurve:new D(1,1,3,7),toneDeltaPair:e=>new z(c.primaryContainer,c.primary,15,"nearer",!1)});c.onPrimaryContainer=C.fromPalette({name:"on_primary_container",palette:e=>e.primaryPalette,tone:e=>yt(e)?C.foregroundTone(c.primaryContainer.tone(e),4.5):R(e)?e.isDark?0:100:e.isDark?90:10,background:e=>c.primaryContainer,contrastCurve:new D(4.5,7,11,21)});c.inversePrimary=C.fromPalette({name:"inverse_primary",palette:e=>e.primaryPalette,tone:e=>e.isDark?40:80,background:e=>c.inverseSurface,contrastCurve:new D(3,4.5,7,11)});c.secondary=C.fromPalette({name:"secondary",palette:e=>e.secondaryPalette,tone:e=>e.isDark?80:40,isBackground:!0,background:e=>c.highestSurface(e),contrastCurve:new D(3,4.5,7,11),toneDeltaPair:e=>new z(c.secondaryContainer,c.secondary,15,"nearer",!1)});c.onSecondary=C.fromPalette({name:"on_secondary",palette:e=>e.secondaryPalette,tone:e=>R(e)?e.isDark?10:100:e.isDark?20:100,background:e=>c.secondary,contrastCurve:new D(4.5,7,11,21)});c.secondaryContainer=C.fromPalette({name:"secondary_container",palette:e=>e.secondaryPalette,tone:e=>{let t=e.isDark?30:90;if(R(e))return e.isDark?30:85;if(!yt(e))return t;let r=He(e.secondaryPalette.hue,e.secondaryPalette.chroma,t,!e.isDark);return r=Zt(e.secondaryPalette.getHct(r),e),r},isBackground:!0,background:e=>c.highestSurface(e),contrastCurve:new D(1,1,3,7),toneDeltaPair:e=>new z(c.secondaryContainer,c.secondary,15,"nearer",!1)});c.onSecondaryContainer=C.fromPalette({name:"on_secondary_container",palette:e=>e.secondaryPalette,tone:e=>yt(e)?C.foregroundTone(c.secondaryContainer.tone(e),4.5):e.isDark?90:10,background:e=>c.secondaryContainer,contrastCurve:new D(4.5,7,11,21)});c.tertiary=C.fromPalette({name:"tertiary",palette:e=>e.tertiaryPalette,tone:e=>R(e)?e.isDark?90:25:e.isDark?80:40,isBackground:!0,background:e=>c.highestSurface(e),contrastCurve:new D(3,4.5,7,11),toneDeltaPair:e=>new z(c.tertiaryContainer,c.tertiary,15,"nearer",!1)});c.onTertiary=C.fromPalette({name:"on_tertiary",palette:e=>e.tertiaryPalette,tone:e=>R(e)?e.isDark?10:90:e.isDark?20:100,background:e=>c.tertiary,contrastCurve:new D(4.5,7,11,21)});c.tertiaryContainer=C.fromPalette({name:"tertiary_container",palette:e=>e.tertiaryPalette,tone:e=>{if(R(e))return e.isDark?60:49;if(!yt(e))return e.isDark?30:90;let t=Zt(e.tertiaryPalette.getHct(e.sourceColorHct.tone),e),r=e.tertiaryPalette.getHct(t);return tt.fixIfDisliked(r).tone},isBackground:!0,background:e=>c.highestSurface(e),contrastCurve:new D(1,1,3,7),toneDeltaPair:e=>new z(c.tertiaryContainer,c.tertiary,15,"nearer",!1)});c.onTertiaryContainer=C.fromPalette({name:"on_tertiary_container",palette:e=>e.tertiaryPalette,tone:e=>R(e)?e.isDark?0:100:yt(e)?C.foregroundTone(c.tertiaryContainer.tone(e),4.5):e.isDark?90:10,background:e=>c.tertiaryContainer,contrastCurve:new D(4.5,7,11,21)});c.error=C.fromPalette({name:"error",palette:e=>e.errorPalette,tone:e=>e.isDark?80:40,isBackground:!0,background:e=>c.highestSurface(e),contrastCurve:new D(3,4.5,7,11),toneDeltaPair:e=>new z(c.errorContainer,c.error,15,"nearer",!1)});c.onError=C.fromPalette({name:"on_error",palette:e=>e.errorPalette,tone:e=>e.isDark?20:100,background:e=>c.error,contrastCurve:new D(4.5,7,11,21)});c.errorContainer=C.fromPalette({name:"error_container",palette:e=>e.errorPalette,tone:e=>e.isDark?30:90,isBackground:!0,background:e=>c.highestSurface(e),contrastCurve:new D(1,1,3,7),toneDeltaPair:e=>new z(c.errorContainer,c.error,15,"nearer",!1)});c.onErrorContainer=C.fromPalette({name:"on_error_container",palette:e=>e.errorPalette,tone:e=>e.isDark?90:10,background:e=>c.errorContainer,contrastCurve:new D(4.5,7,11,21)});c.primaryFixed=C.fromPalette({name:"primary_fixed",palette:e=>e.primaryPalette,tone:e=>R(e)?40:90,isBackground:!0,background:e=>c.highestSurface(e),contrastCurve:new D(1,1,3,7),toneDeltaPair:e=>new z(c.primaryFixed,c.primaryFixedDim,10,"lighter",!0)});c.primaryFixedDim=C.fromPalette({name:"primary_fixed_dim",palette:e=>e.primaryPalette,tone:e=>R(e)?30:80,isBackground:!0,background:e=>c.highestSurface(e),contrastCurve:new D(1,1,3,7),toneDeltaPair:e=>new z(c.primaryFixed,c.primaryFixedDim,10,"lighter",!0)});c.onPrimaryFixed=C.fromPalette({name:"on_primary_fixed",palette:e=>e.primaryPalette,tone:e=>R(e)?100:10,background:e=>c.primaryFixedDim,secondBackground:e=>c.primaryFixed,contrastCurve:new D(4.5,7,11,21)});c.onPrimaryFixedVariant=C.fromPalette({name:"on_primary_fixed_variant",palette:e=>e.primaryPalette,tone:e=>R(e)?90:30,background:e=>c.primaryFixedDim,secondBackground:e=>c.primaryFixed,contrastCurve:new D(3,4.5,7,11)});c.secondaryFixed=C.fromPalette({name:"secondary_fixed",palette:e=>e.secondaryPalette,tone:e=>R(e)?80:90,isBackground:!0,background:e=>c.highestSurface(e),contrastCurve:new D(1,1,3,7),toneDeltaPair:e=>new z(c.secondaryFixed,c.secondaryFixedDim,10,"lighter",!0)});c.secondaryFixedDim=C.fromPalette({name:"secondary_fixed_dim",palette:e=>e.secondaryPalette,tone:e=>R(e)?70:80,isBackground:!0,background:e=>c.highestSurface(e),contrastCurve:new D(1,1,3,7),toneDeltaPair:e=>new z(c.secondaryFixed,c.secondaryFixedDim,10,"lighter",!0)});c.onSecondaryFixed=C.fromPalette({name:"on_secondary_fixed",palette:e=>e.secondaryPalette,tone:e=>10,background:e=>c.secondaryFixedDim,secondBackground:e=>c.secondaryFixed,contrastCurve:new D(4.5,7,11,21)});c.onSecondaryFixedVariant=C.fromPalette({name:"on_secondary_fixed_variant",palette:e=>e.secondaryPalette,tone:e=>R(e)?25:30,background:e=>c.secondaryFixedDim,secondBackground:e=>c.secondaryFixed,contrastCurve:new D(3,4.5,7,11)});c.tertiaryFixed=C.fromPalette({name:"tertiary_fixed",palette:e=>e.tertiaryPalette,tone:e=>R(e)?40:90,isBackground:!0,background:e=>c.highestSurface(e),contrastCurve:new D(1,1,3,7),toneDeltaPair:e=>new z(c.tertiaryFixed,c.tertiaryFixedDim,10,"lighter",!0)});c.tertiaryFixedDim=C.fromPalette({name:"tertiary_fixed_dim",palette:e=>e.tertiaryPalette,tone:e=>R(e)?30:80,isBackground:!0,background:e=>c.highestSurface(e),contrastCurve:new D(1,1,3,7),toneDeltaPair:e=>new z(c.tertiaryFixed,c.tertiaryFixedDim,10,"lighter",!0)});c.onTertiaryFixed=C.fromPalette({name:"on_tertiary_fixed",palette:e=>e.tertiaryPalette,tone:e=>R(e)?100:10,background:e=>c.tertiaryFixedDim,secondBackground:e=>c.tertiaryFixed,contrastCurve:new D(4.5,7,11,21)});c.onTertiaryFixedVariant=C.fromPalette({name:"on_tertiary_fixed_variant",palette:e=>e.tertiaryPalette,tone:e=>R(e)?90:30,background:e=>c.tertiaryFixedDim,secondBackground:e=>c.tertiaryFixed,contrastCurve:new D(3,4.5,7,11)});var A=class e{static fromInt(t){let r=F.fromInt(t);return e.fromHct(r)}static fromHct(t){return new e(t.hue,t.chroma,t)}static fromHueAndChroma(t,r){return new e(t,r,e.createKeyColor(t,r))}constructor(t,r,n){this.hue=t,this.chroma=r,this.keyColor=n,this.cache=new Map}static createKeyColor(t,r){let o=F.from(t,r,50),a=Math.abs(o.chroma-r);for(let s=1;s<50;s+=1){if(Math.round(r)===Math.round(o.chroma))return o;let i=F.from(t,r,50+s),l=Math.abs(i.chroma-r);l<a&&(a=l,o=i);let u=F.from(t,r,50-s),h=Math.abs(u.chroma-r);h<a&&(a=h,o=u)}return o}tone(t){let r=this.cache.get(t);return r===void 0&&(r=F.from(this.hue,this.chroma,t).toInt(),this.cache.set(t,r)),r}getHct(t){return F.fromInt(this.tone(t))}};var G=class e{static of(t){return new e(t,!1)}static contentOf(t){return new e(t,!0)}static fromColors(t){return e.createPaletteFromColors(!1,t)}static contentFromColors(t){return e.createPaletteFromColors(!0,t)}static createPaletteFromColors(t,r){let n=new e(r.primary,t);if(r.secondary){let o=new e(r.secondary,t);n.a2=o.a1}if(r.tertiary){let o=new e(r.tertiary,t);n.a3=o.a1}if(r.error){let o=new e(r.error,t);n.error=o.a1}if(r.neutral){let o=new e(r.neutral,t);n.n1=o.n1}if(r.neutralVariant){let o=new e(r.neutralVariant,t);n.n2=o.n2}return n}constructor(t,r){let n=F.fromInt(t),o=n.hue,a=n.chroma;r?(this.a1=A.fromHueAndChroma(o,a),this.a2=A.fromHueAndChroma(o,a/3),this.a3=A.fromHueAndChroma(o+60,a/2),this.n1=A.fromHueAndChroma(o,Math.min(a/12,4)),this.n2=A.fromHueAndChroma(o,Math.min(a/6,8))):(this.a1=A.fromHueAndChroma(o,Math.max(48,a)),this.a2=A.fromHueAndChroma(o,16),this.a3=A.fromHueAndChroma(o+60,24),this.n1=A.fromHueAndChroma(o,4),this.n2=A.fromHueAndChroma(o,8)),this.error=A.fromHueAndChroma(25,84)}};var Vt=class{fromInt(t){return kt(t)}toInt(t){return Jt(t[0],t[1],t[2])}distance(t,r){let n=t[0]-r[0],o=t[1]-r[1],a=t[2]-r[2];return n*n+o*o+a*a}};var Ee=10,Le=3,Dt=class{static quantize(t,r,n){let o=new Map,a=new Array,s=new Array,i=new Vt,l=0;for(let g=0;g<t.length;g++){let P=t[g],T=o.get(P);T===void 0?(l++,a.push(i.fromInt(P)),s.push(P),o.set(P,1)):o.set(P,T+1)}let u=new Array;for(let g=0;g<l;g++){let P=s[g],T=o.get(P);T!==void 0&&(u[g]=T)}let h=Math.min(n,l);r.length>0&&(h=Math.min(h,r.length));let p=new Array;for(let g=0;g<r.length;g++)p.push(i.fromInt(r[g]));let y=h-p.length;if(r.length===0&&y>0)for(let g=0;g<y;g++){let P=Math.random()*100,T=Math.random()*(100- -100+1)+-100,I=Math.random()*(100- -100+1)+-100;p.push(new Array(P,T,I))}let m=new Array;for(let g=0;g<l;g++)m.push(Math.floor(Math.random()*h));let d=new Array;for(let g=0;g<h;g++){d.push(new Array);for(let P=0;P<h;P++)d[g].push(0)}let f=new Array;for(let g=0;g<h;g++){f.push(new Array);for(let P=0;P<h;P++)f[g].push(new Qt)}let b=new Array;for(let g=0;g<h;g++)b.push(0);for(let g=0;g<Ee;g++){for(let x=0;x<h;x++){for(let M=x+1;M<h;M++){let B=i.distance(p[x],p[M]);f[M][x].distance=B,f[M][x].index=x,f[x][M].distance=B,f[x][M].index=M}f[x].sort();for(let M=0;M<h;M++)d[x][M]=f[x][M].index}let P=0;for(let x=0;x<l;x++){let M=a[x],B=m[x],w=p[B],q=i.distance(M,w),W=q,Y=-1;for(let E=0;E<h;E++){if(f[B][E].distance>=4*q)continue;let ct=i.distance(M,p[E]);ct<W&&(W=ct,Y=E)}Y!==-1&&Math.abs(Math.sqrt(W)-Math.sqrt(q))>Le&&(P++,m[x]=Y)}if(P===0&&g!==0)break;let T=new Array(h).fill(0),I=new Array(h).fill(0),S=new Array(h).fill(0);for(let x=0;x<h;x++)b[x]=0;for(let x=0;x<l;x++){let M=m[x],B=a[x],w=u[x];b[M]+=w,T[M]+=B[0]*w,I[M]+=B[1]*w,S[M]+=B[2]*w}for(let x=0;x<h;x++){let M=b[x];if(M===0){p[x]=[0,0,0];continue}let B=T[x]/M,w=I[x]/M,q=S[x]/M;p[x]=[B,w,q]}}let k=new Map;for(let g=0;g<h;g++){let P=b[g];if(P===0)continue;let T=i.toInt(p[g]);k.has(T)||k.set(T,P)}return k}},Qt=class{constructor(){this.distance=-1,this.index=-1}};var Tt=class{static quantize(t){var n;let r=new Map;for(let o=0;o<t.length;o++){let a=t[o];Mt(a)<255||r.set(a,((n=r.get(a))!=null?n:0)+1)}return r}};var Nt=5,Z=33,wt=35937,j={RED:"red",GREEN:"green",BLUE:"blue"},Bt=class{constructor(t=[],r=[],n=[],o=[],a=[],s=[]){this.weights=t,this.momentsR=r,this.momentsG=n,this.momentsB=o,this.moments=a,this.cubes=s}quantize(t,r){this.constructHistogram(t),this.computeMoments();let n=this.createBoxes(r);return this.createResult(n.resultCount)}constructHistogram(t){var n;this.weights=Array.from({length:wt}).fill(0),this.momentsR=Array.from({length:wt}).fill(0),this.momentsG=Array.from({length:wt}).fill(0),this.momentsB=Array.from({length:wt}).fill(0),this.moments=Array.from({length:wt}).fill(0);let r=Tt.quantize(t);for(let[o,a]of r.entries()){let s=at(o),i=st(o),l=it(o),u=8-Nt,h=(s>>u)+1,p=(i>>u)+1,y=(l>>u)+1,m=this.getIndex(h,p,y);this.weights[m]=((n=this.weights[m])!=null?n:0)+a,this.momentsR[m]+=a*s,this.momentsG[m]+=a*i,this.momentsB[m]+=a*l,this.moments[m]+=a*(s*s+i*i+l*l)}}computeMoments(){for(let t=1;t<Z;t++){let r=Array.from({length:Z}).fill(0),n=Array.from({length:Z}).fill(0),o=Array.from({length:Z}).fill(0),a=Array.from({length:Z}).fill(0),s=Array.from({length:Z}).fill(0);for(let i=1;i<Z;i++){let l=0,u=0,h=0,p=0,y=0;for(let m=1;m<Z;m++){let d=this.getIndex(t,i,m);l+=this.weights[d],u+=this.momentsR[d],h+=this.momentsG[d],p+=this.momentsB[d],y+=this.moments[d],r[m]+=l,n[m]+=u,o[m]+=h,a[m]+=p,s[m]+=y;let f=this.getIndex(t-1,i,m);this.weights[d]=this.weights[f]+r[m],this.momentsR[d]=this.momentsR[f]+n[m],this.momentsG[d]=this.momentsG[f]+o[m],this.momentsB[d]=this.momentsB[f]+a[m],this.moments[d]=this.moments[f]+s[m]}}}}createBoxes(t){this.cubes=Array.from({length:t}).fill(0).map(()=>new te);let r=Array.from({length:t}).fill(0);this.cubes[0].r0=0,this.cubes[0].g0=0,this.cubes[0].b0=0,this.cubes[0].r1=Z-1,this.cubes[0].g1=Z-1,this.cubes[0].b1=Z-1;let n=t,o=0;for(let a=1;a<t;a++){this.cut(this.cubes[o],this.cubes[a])?(r[o]=this.cubes[o].vol>1?this.variance(this.cubes[o]):0,r[a]=this.cubes[a].vol>1?this.variance(this.cubes[a]):0):(r[o]=0,a--),o=0;let s=r[0];for(let i=1;i<=a;i++)r[i]>s&&(s=r[i],o=i);if(s<=0){n=a+1;break}}return new ee(t,n)}createResult(t){let r=[];for(let n=0;n<t;++n){let o=this.cubes[n],a=this.volume(o,this.weights);if(a>0){let s=Math.round(this.volume(o,this.momentsR)/a),i=Math.round(this.volume(o,this.momentsG)/a),l=Math.round(this.volume(o,this.momentsB)/a),u=255<<24|(s&255)<<16|(i&255)<<8|l&255;r.push(u)}}return r}variance(t){let r=this.volume(t,this.momentsR),n=this.volume(t,this.momentsG),o=this.volume(t,this.momentsB),a=this.moments[this.getIndex(t.r1,t.g1,t.b1)]-this.moments[this.getIndex(t.r1,t.g1,t.b0)]-this.moments[this.getIndex(t.r1,t.g0,t.b1)]+this.moments[this.getIndex(t.r1,t.g0,t.b0)]-this.moments[this.getIndex(t.r0,t.g1,t.b1)]+this.moments[this.getIndex(t.r0,t.g1,t.b0)]+this.moments[this.getIndex(t.r0,t.g0,t.b1)]-this.moments[this.getIndex(t.r0,t.g0,t.b0)],s=r*r+n*n+o*o,i=this.volume(t,this.weights);return a-s/i}cut(t,r){let n=this.volume(t,this.momentsR),o=this.volume(t,this.momentsG),a=this.volume(t,this.momentsB),s=this.volume(t,this.weights),i=this.maximize(t,j.RED,t.r0+1,t.r1,n,o,a,s),l=this.maximize(t,j.GREEN,t.g0+1,t.g1,n,o,a,s),u=this.maximize(t,j.BLUE,t.b0+1,t.b1,n,o,a,s),h,p=i.maximum,y=l.maximum,m=u.maximum;if(p>=y&&p>=m){if(i.cutLocation<0)return!1;h=j.RED}else y>=p&&y>=m?h=j.GREEN:h=j.BLUE;switch(r.r1=t.r1,r.g1=t.g1,r.b1=t.b1,h){case j.RED:t.r1=i.cutLocation,r.r0=t.r1,r.g0=t.g0,r.b0=t.b0;break;case j.GREEN:t.g1=l.cutLocation,r.r0=t.r0,r.g0=t.g1,r.b0=t.b0;break;case j.BLUE:t.b1=u.cutLocation,r.r0=t.r0,r.g0=t.g0,r.b0=t.b1;break;default:throw new Error("unexpected direction "+h)}return t.vol=(t.r1-t.r0)*(t.g1-t.g0)*(t.b1-t.b0),r.vol=(r.r1-r.r0)*(r.g1-r.g0)*(r.b1-r.b0),!0}maximize(t,r,n,o,a,s,i,l){let u=this.bottom(t,r,this.momentsR),h=this.bottom(t,r,this.momentsG),p=this.bottom(t,r,this.momentsB),y=this.bottom(t,r,this.weights),m=0,d=-1,f=0,b=0,k=0,g=0;for(let P=n;P<o;P++){if(f=u+this.top(t,r,P,this.momentsR),b=h+this.top(t,r,P,this.momentsG),k=p+this.top(t,r,P,this.momentsB),g=y+this.top(t,r,P,this.weights),g===0)continue;let T=(f*f+b*b+k*k)*1,I=g*1,S=T/I;f=a-f,b=s-b,k=i-k,g=l-g,g!==0&&(T=(f*f+b*b+k*k)*1,I=g*1,S+=T/I,S>m&&(m=S,d=P))}return new re(d,m)}volume(t,r){return r[this.getIndex(t.r1,t.g1,t.b1)]-r[this.getIndex(t.r1,t.g1,t.b0)]-r[this.getIndex(t.r1,t.g0,t.b1)]+r[this.getIndex(t.r1,t.g0,t.b0)]-r[this.getIndex(t.r0,t.g1,t.b1)]+r[this.getIndex(t.r0,t.g1,t.b0)]+r[this.getIndex(t.r0,t.g0,t.b1)]-r[this.getIndex(t.r0,t.g0,t.b0)]}bottom(t,r,n){switch(r){case j.RED:return-n[this.getIndex(t.r0,t.g1,t.b1)]+n[this.getIndex(t.r0,t.g1,t.b0)]+n[this.getIndex(t.r0,t.g0,t.b1)]-n[this.getIndex(t.r0,t.g0,t.b0)];case j.GREEN:return-n[this.getIndex(t.r1,t.g0,t.b1)]+n[this.getIndex(t.r1,t.g0,t.b0)]+n[this.getIndex(t.r0,t.g0,t.b1)]-n[this.getIndex(t.r0,t.g0,t.b0)];case j.BLUE:return-n[this.getIndex(t.r1,t.g1,t.b0)]+n[this.getIndex(t.r1,t.g0,t.b0)]+n[this.getIndex(t.r0,t.g1,t.b0)]-n[this.getIndex(t.r0,t.g0,t.b0)];default:throw new Error("unexpected direction $direction")}}top(t,r,n,o){switch(r){case j.RED:return o[this.getIndex(n,t.g1,t.b1)]-o[this.getIndex(n,t.g1,t.b0)]-o[this.getIndex(n,t.g0,t.b1)]+o[this.getIndex(n,t.g0,t.b0)];case j.GREEN:return o[this.getIndex(t.r1,n,t.b1)]-o[this.getIndex(t.r1,n,t.b0)]-o[this.getIndex(t.r0,n,t.b1)]+o[this.getIndex(t.r0,n,t.b0)];case j.BLUE:return o[this.getIndex(t.r1,t.g1,n)]-o[this.getIndex(t.r1,t.g0,n)]-o[this.getIndex(t.r0,t.g1,n)]+o[this.getIndex(t.r0,t.g0,n)];default:throw new Error("unexpected direction $direction")}}getIndex(t,r,n){return(t<<Nt*2)+(t<<Nt+1)+t+(r<<Nt)+r+n}},te=class{constructor(t=0,r=0,n=0,o=0,a=0,s=0,i=0){this.r0=t,this.r1=r,this.g0=n,this.g1=o,this.b0=a,this.b1=s,this.vol=i}},ee=class{constructor(t,r){this.requestedCount=t,this.resultCount=r}},re=class{constructor(t,r){this.cutLocation=t,this.maximum=r}};var Ft=class{static quantize(t,r){let o=new Bt().quantize(t,r);return Dt.quantize(t,o,r)}};var H=class{constructor(t){this.sourceColorArgb=t.sourceColorArgb,this.variant=t.variant,this.contrastLevel=t.contrastLevel,this.isDark=t.isDark,this.sourceColorHct=F.fromInt(t.sourceColorArgb),this.primaryPalette=t.primaryPalette,this.secondaryPalette=t.secondaryPalette,this.tertiaryPalette=t.tertiaryPalette,this.neutralPalette=t.neutralPalette,this.neutralVariantPalette=t.neutralVariantPalette,this.errorPalette=A.fromHueAndChroma(25,84)}static getRotatedHue(t,r,n){let o=t.hue;if(r.length!==n.length)throw new Error(`mismatch between hue length ${r.length} & rotations ${n.length}`);if(n.length===1)return v(t.hue+n[0]);let a=r.length;for(let s=0;s<=a-2;s++){let i=r[s],l=r[s+1];if(i<o&&o<l)return v(o+n[s])}return o}};var Pt=class e{get primary(){return this.props.primary}get onPrimary(){return this.props.onPrimary}get primaryContainer(){return this.props.primaryContainer}get onPrimaryContainer(){return this.props.onPrimaryContainer}get secondary(){return this.props.secondary}get onSecondary(){return this.props.onSecondary}get secondaryContainer(){return this.props.secondaryContainer}get onSecondaryContainer(){return this.props.onSecondaryContainer}get tertiary(){return this.props.tertiary}get onTertiary(){return this.props.onTertiary}get tertiaryContainer(){return this.props.tertiaryContainer}get onTertiaryContainer(){return this.props.onTertiaryContainer}get error(){return this.props.error}get onError(){return this.props.onError}get errorContainer(){return this.props.errorContainer}get onErrorContainer(){return this.props.onErrorContainer}get background(){return this.props.background}get onBackground(){return this.props.onBackground}get surface(){return this.props.surface}get onSurface(){return this.props.onSurface}get surfaceVariant(){return this.props.surfaceVariant}get onSurfaceVariant(){return this.props.onSurfaceVariant}get outline(){return this.props.outline}get outlineVariant(){return this.props.outlineVariant}get shadow(){return this.props.shadow}get scrim(){return this.props.scrim}get inverseSurface(){return this.props.inverseSurface}get inverseOnSurface(){return this.props.inverseOnSurface}get inversePrimary(){return this.props.inversePrimary}static light(t){return e.lightFromCorePalette(G.of(t))}static dark(t){return e.darkFromCorePalette(G.of(t))}static lightContent(t){return e.lightFromCorePalette(G.contentOf(t))}static darkContent(t){return e.darkFromCorePalette(G.contentOf(t))}static lightFromCorePalette(t){return new e({primary:t.a1.tone(40),onPrimary:t.a1.tone(100),primaryContainer:t.a1.tone(90),onPrimaryContainer:t.a1.tone(10),secondary:t.a2.tone(40),onSecondary:t.a2.tone(100),secondaryContainer:t.a2.tone(90),onSecondaryContainer:t.a2.tone(10),tertiary:t.a3.tone(40),onTertiary:t.a3.tone(100),tertiaryContainer:t.a3.tone(90),onTertiaryContainer:t.a3.tone(10),error:t.error.tone(40),onError:t.error.tone(100),errorContainer:t.error.tone(90),onErrorContainer:t.error.tone(10),background:t.n1.tone(99),onBackground:t.n1.tone(10),surface:t.n1.tone(99),onSurface:t.n1.tone(10),surfaceVariant:t.n2.tone(90),onSurfaceVariant:t.n2.tone(30),outline:t.n2.tone(50),outlineVariant:t.n2.tone(80),shadow:t.n1.tone(0),scrim:t.n1.tone(0),inverseSurface:t.n1.tone(20),inverseOnSurface:t.n1.tone(95),inversePrimary:t.a1.tone(80)})}static darkFromCorePalette(t){return new e({primary:t.a1.tone(80),onPrimary:t.a1.tone(20),primaryContainer:t.a1.tone(30),onPrimaryContainer:t.a1.tone(90),secondary:t.a2.tone(80),onSecondary:t.a2.tone(20),secondaryContainer:t.a2.tone(30),onSecondaryContainer:t.a2.tone(90),tertiary:t.a3.tone(80),onTertiary:t.a3.tone(20),tertiaryContainer:t.a3.tone(30),onTertiaryContainer:t.a3.tone(90),error:t.error.tone(80),onError:t.error.tone(20),errorContainer:t.error.tone(30),onErrorContainer:t.error.tone(80),background:t.n1.tone(10),onBackground:t.n1.tone(90),surface:t.n1.tone(10),onSurface:t.n1.tone(90),surfaceVariant:t.n2.tone(30),onSurfaceVariant:t.n2.tone(80),outline:t.n2.tone(60),outlineVariant:t.n2.tone(30),shadow:t.n1.tone(0),scrim:t.n1.tone(0),inverseSurface:t.n1.tone(90),inverseOnSurface:t.n1.tone(20),inversePrimary:t.a1.tone(40)})}constructor(t){this.props=t}toJSON(){return lt({},this.props)}};var ne=class e{get colorAccentPrimary(){return this.props.colorAccentPrimary}get colorAccentPrimaryVariant(){return this.props.colorAccentPrimaryVariant}get colorAccentSecondary(){return this.props.colorAccentSecondary}get colorAccentSecondaryVariant(){return this.props.colorAccentSecondaryVariant}get colorAccentTertiary(){return this.props.colorAccentTertiary}get colorAccentTertiaryVariant(){return this.props.colorAccentTertiaryVariant}get textColorPrimary(){return this.props.textColorPrimary}get textColorSecondary(){return this.props.textColorSecondary}get textColorTertiary(){return this.props.textColorTertiary}get textColorPrimaryInverse(){return this.props.textColorPrimaryInverse}get textColorSecondaryInverse(){return this.props.textColorSecondaryInverse}get textColorTertiaryInverse(){return this.props.textColorTertiaryInverse}get colorBackground(){return this.props.colorBackground}get colorBackgroundFloating(){return this.props.colorBackgroundFloating}get colorSurface(){return this.props.colorSurface}get colorSurfaceVariant(){return this.props.colorSurfaceVariant}get colorSurfaceHighlight(){return this.props.colorSurfaceHighlight}get surfaceHeader(){return this.props.surfaceHeader}get underSurface(){return this.props.underSurface}get offState(){return this.props.offState}get accentSurface(){return this.props.accentSurface}get textPrimaryOnAccent(){return this.props.textPrimaryOnAccent}get textSecondaryOnAccent(){return this.props.textSecondaryOnAccent}get volumeBackground(){return this.props.volumeBackground}get scrim(){return this.props.scrim}static light(t){let r=G.of(t);return e.lightFromCorePalette(r)}static dark(t){let r=G.of(t);return e.darkFromCorePalette(r)}static lightContent(t){let r=G.contentOf(t);return e.lightFromCorePalette(r)}static darkContent(t){let r=G.contentOf(t);return e.darkFromCorePalette(r)}static lightFromCorePalette(t){return new e({colorAccentPrimary:t.a1.tone(90),colorAccentPrimaryVariant:t.a1.tone(40),colorAccentSecondary:t.a2.tone(90),colorAccentSecondaryVariant:t.a2.tone(40),colorAccentTertiary:t.a3.tone(90),colorAccentTertiaryVariant:t.a3.tone(40),textColorPrimary:t.n1.tone(10),textColorSecondary:t.n2.tone(30),textColorTertiary:t.n2.tone(50),textColorPrimaryInverse:t.n1.tone(95),textColorSecondaryInverse:t.n1.tone(80),textColorTertiaryInverse:t.n1.tone(60),colorBackground:t.n1.tone(95),colorBackgroundFloating:t.n1.tone(98),colorSurface:t.n1.tone(98),colorSurfaceVariant:t.n1.tone(90),colorSurfaceHighlight:t.n1.tone(100),surfaceHeader:t.n1.tone(90),underSurface:t.n1.tone(0),offState:t.n1.tone(20),accentSurface:t.a2.tone(95),textPrimaryOnAccent:t.n1.tone(10),textSecondaryOnAccent:t.n2.tone(30),volumeBackground:t.n1.tone(25),scrim:t.n1.tone(80)})}static darkFromCorePalette(t){return new e({colorAccentPrimary:t.a1.tone(90),colorAccentPrimaryVariant:t.a1.tone(70),colorAccentSecondary:t.a2.tone(90),colorAccentSecondaryVariant:t.a2.tone(70),colorAccentTertiary:t.a3.tone(90),colorAccentTertiaryVariant:t.a3.tone(70),textColorPrimary:t.n1.tone(95),textColorSecondary:t.n2.tone(80),textColorTertiary:t.n2.tone(60),textColorPrimaryInverse:t.n1.tone(10),textColorSecondaryInverse:t.n1.tone(30),textColorTertiaryInverse:t.n1.tone(50),colorBackground:t.n1.tone(10),colorBackgroundFloating:t.n1.tone(10),colorSurface:t.n1.tone(20),colorSurfaceVariant:t.n1.tone(30),colorSurfaceHighlight:t.n1.tone(35),surfaceHeader:t.n1.tone(30),underSurface:t.n1.tone(0),offState:t.n1.tone(20),accentSurface:t.a2.tone(95),textPrimaryOnAccent:t.n1.tone(10),textSecondaryOnAccent:t.n2.tone(30),volumeBackground:t.n1.tone(25),scrim:t.n1.tone(80)})}constructor(t){this.props=t}toJSON(){return lt({},this.props)}};var mt=class e{constructor(t){this.input=t,this.hctsByTempCache=[],this.hctsByHueCache=[],this.tempsByHctCache=new Map,this.inputRelativeTemperatureCache=-1,this.complementCache=null}get hctsByTemp(){if(this.hctsByTempCache.length>0)return this.hctsByTempCache;let t=this.hctsByHue.concat([this.input]),r=this.tempsByHct;return t.sort((n,o)=>r.get(n)-r.get(o)),this.hctsByTempCache=t,t}get warmest(){return this.hctsByTemp[this.hctsByTemp.length-1]}get coldest(){return this.hctsByTemp[0]}analogous(t=5,r=12){let n=Math.round(this.input.hue),o=this.hctsByHue[n],a=this.relativeTemperature(o),s=[o],i=0;for(let d=0;d<360;d++){let f=ut(n+d),b=this.hctsByHue[f],k=this.relativeTemperature(b),g=Math.abs(k-a);a=k,i+=g}let l=1,u=i/r,h=0;for(a=this.relativeTemperature(o);s.length<r;){let d=ut(n+l),f=this.hctsByHue[d],b=this.relativeTemperature(f),k=Math.abs(b-a);h+=k;let g=s.length*u,P=h>=g,T=1;for(;P&&s.length<r;){s.push(f);let I=(s.length+T)*u;P=h>=I,T++}if(a=b,l++,l>360){for(;s.length<r;)s.push(f);break}}let p=[this.input],y=Math.floor((t-1)/2);for(let d=1;d<y+1;d++){let f=0-d;for(;f<0;)f=s.length+f;f>=s.length&&(f=f%s.length),p.splice(0,0,s[f])}let m=t-y-1;for(let d=1;d<m+1;d++){let f=d;for(;f<0;)f=s.length+f;f>=s.length&&(f=f%s.length),p.push(s[f])}return p}get complement(){if(this.complementCache!=null)return this.complementCache;let t=this.coldest.hue,r=this.tempsByHct.get(this.coldest),n=this.warmest.hue,a=this.tempsByHct.get(this.warmest)-r,s=e.isBetween(this.input.hue,t,n),i=s?n:t,l=s?t:n,u=1,h=1e3,p=this.hctsByHue[Math.round(this.input.hue)],y=1-this.inputRelativeTemperature;for(let m=0;m<=360;m+=1){let d=v(i+u*m);if(!e.isBetween(d,i,l))continue;let f=this.hctsByHue[Math.round(d)],b=(this.tempsByHct.get(f)-r)/a,k=Math.abs(y-b);k<h&&(h=k,p=f)}return this.complementCache=p,this.complementCache}relativeTemperature(t){let r=this.tempsByHct.get(this.warmest)-this.tempsByHct.get(this.coldest),n=this.tempsByHct.get(t)-this.tempsByHct.get(this.coldest);return r===0?.5:n/r}get inputRelativeTemperature(){return this.inputRelativeTemperatureCache>=0?this.inputRelativeTemperatureCache:(this.inputRelativeTemperatureCache=this.relativeTemperature(this.input),this.inputRelativeTemperatureCache)}get tempsByHct(){if(this.tempsByHctCache.size>0)return this.tempsByHctCache;let t=this.hctsByHue.concat([this.input]),r=new Map;for(let n of t)r.set(n,e.rawTemperature(n));return this.tempsByHctCache=r,r}get hctsByHue(){if(this.hctsByHueCache.length>0)return this.hctsByHueCache;let t=[];for(let r=0;r<=360;r+=1){let n=F.from(r,this.input.chroma,this.input.tone);t.push(n)}return this.hctsByHueCache=t,this.hctsByHueCache}static isBetween(t,r,n){return r<n?r<=t&&t<=n:r<=t||t<=n}static rawTemperature(t){let r=kt(t.toInt()),n=v(Math.atan2(r[2],r[1])*180/Math.PI),o=Math.sqrt(r[1]*r[1]+r[2]*r[2]);return-.5+.02*Math.pow(o,1.07)*Math.cos(v(n-50)*Math.PI/180)}};var oe=class extends H{constructor(t,r,n){super({sourceColorArgb:t.toInt(),variant:O.CONTENT,contrastLevel:n,isDark:r,primaryPalette:A.fromHueAndChroma(t.hue,t.chroma),secondaryPalette:A.fromHueAndChroma(t.hue,Math.max(t.chroma-32,t.chroma*.5)),tertiaryPalette:A.fromInt(tt.fixIfDisliked(new mt(t).analogous(3,6)[2]).toInt()),neutralPalette:A.fromHueAndChroma(t.hue,t.chroma/8),neutralVariantPalette:A.fromHueAndChroma(t.hue,t.chroma/8+4)})}};var xt=class e extends H{constructor(t,r,n){super({sourceColorArgb:t.toInt(),variant:O.EXPRESSIVE,contrastLevel:n,isDark:r,primaryPalette:A.fromHueAndChroma(v(t.hue+240),40),secondaryPalette:A.fromHueAndChroma(H.getRotatedHue(t,e.hues,e.secondaryRotations),24),tertiaryPalette:A.fromHueAndChroma(H.getRotatedHue(t,e.hues,e.tertiaryRotations),32),neutralPalette:A.fromHueAndChroma(t.hue+15,8),neutralVariantPalette:A.fromHueAndChroma(t.hue+15,12)})}};xt.hues=[0,21,51,121,151,191,271,321,360];xt.secondaryRotations=[45,95,45,20,45,90,45,45,45];xt.tertiaryRotations=[120,120,20,45,20,15,20,120,120];var ae=class extends H{constructor(t,r,n){super({sourceColorArgb:t.toInt(),variant:O.FIDELITY,contrastLevel:n,isDark:r,primaryPalette:A.fromHueAndChroma(t.hue,t.chroma),secondaryPalette:A.fromHueAndChroma(t.hue,Math.max(t.chroma-32,t.chroma*.5)),tertiaryPalette:A.fromInt(tt.fixIfDisliked(new mt(t).complement).toInt()),neutralPalette:A.fromHueAndChroma(t.hue,t.chroma/8),neutralVariantPalette:A.fromHueAndChroma(t.hue,t.chroma/8+4)})}};var se=class extends H{constructor(t,r,n){super({sourceColorArgb:t.toInt(),variant:O.MONOCHROME,contrastLevel:n,isDark:r,primaryPalette:A.fromHueAndChroma(t.hue,0),secondaryPalette:A.fromHueAndChroma(t.hue,0),tertiaryPalette:A.fromHueAndChroma(t.hue,0),neutralPalette:A.fromHueAndChroma(t.hue,0),neutralVariantPalette:A.fromHueAndChroma(t.hue,0)})}};var ie=class extends H{constructor(t,r,n){super({sourceColorArgb:t.toInt(),variant:O.NEUTRAL,contrastLevel:n,isDark:r,primaryPalette:A.fromHueAndChroma(t.hue,12),secondaryPalette:A.fromHueAndChroma(t.hue,8),tertiaryPalette:A.fromHueAndChroma(t.hue,16),neutralPalette:A.fromHueAndChroma(t.hue,2),neutralVariantPalette:A.fromHueAndChroma(t.hue,2)})}};var ce=class extends H{constructor(t,r,n){super({sourceColorArgb:t.toInt(),variant:O.TONAL_SPOT,contrastLevel:n,isDark:r,primaryPalette:A.fromHueAndChroma(t.hue,36),secondaryPalette:A.fromHueAndChroma(t.hue,16),tertiaryPalette:A.fromHueAndChroma(v(t.hue+60),24),neutralPalette:A.fromHueAndChroma(t.hue,6),neutralVariantPalette:A.fromHueAndChroma(t.hue,8)})}};var Ct=class e extends H{constructor(t,r,n){super({sourceColorArgb:t.toInt(),variant:O.VIBRANT,contrastLevel:n,isDark:r,primaryPalette:A.fromHueAndChroma(t.hue,200),secondaryPalette:A.fromHueAndChroma(H.getRotatedHue(t,e.hues,e.secondaryRotations),24),tertiaryPalette:A.fromHueAndChroma(H.getRotatedHue(t,e.hues,e.tertiaryRotations),32),neutralPalette:A.fromHueAndChroma(t.hue,10),neutralVariantPalette:A.fromHueAndChroma(t.hue,12)})}};Ct.hues=[0,41,61,101,131,181,251,301,360];Ct.secondaryRotations=[18,15,10,12,15,18,15,12,12];Ct.tertiaryRotations=[35,30,20,25,30,35,30,25,25];var ve={desired:4,fallbackColorARGB:4282549748,filter:!0};function ze(e,t){return e.score>t.score?-1:e.score<t.score?1:0}var Q=class e{constructor(){}static score(t,r){let{desired:n,fallbackColorARGB:o,filter:a}=lt(lt({},ve),r),s=[],i=new Array(360).fill(0),l=0;for(let[m,d]of t.entries()){let f=F.fromInt(m);s.push(f);let b=Math.floor(f.hue);i[b]+=d,l+=d}let u=new Array(360).fill(0);for(let m=0;m<360;m++){let d=i[m]/l;for(let f=m-14;f<m+16;f++){let b=ut(f);u[b]+=d}}let h=new Array;for(let m of s){let d=ut(Math.round(m.hue)),f=u[d];if(a&&(m.chroma<e.CUTOFF_CHROMA||f<=e.CUTOFF_EXCITED_PROPORTION))continue;let b=f*100*e.WEIGHT_PROPORTION,k=m.chroma<e.TARGET_CHROMA?e.WEIGHT_CHROMA_BELOW:e.WEIGHT_CHROMA_ABOVE,g=(m.chroma-e.TARGET_CHROMA)*k,P=b+g;h.push({hct:m,score:P})}h.sort(ze);let p=[];for(let m=90;m>=15;m--){p.length=0;for(let{hct:d}of h)if(p.find(b=>bt(d.hue,b.hue)<m)||p.push(d),p.length>=n)break;if(p.length>=n)break}let y=[];p.length===0&&y.push(o);for(let m of p)y.push(m.toInt());return y}};Q.TARGET_CHROMA=48;Q.WEIGHT_PROPORTION=.7;Q.WEIGHT_CHROMA_ABOVE=.3;Q.WEIGHT_CHROMA_BELOW=.1;Q.CUTOFF_CHROMA=5;Q.CUTOFF_EXCITED_PROPORTION=.01;function Ut(e){let t=at(e),r=st(e),n=it(e),o=[t.toString(16),r.toString(16),n.toString(16)];for(let[a,s]of o.entries())s.length===1&&(o[a]="0"+s);return"#"+o.join("")}function Ve(e){e=e.replace("#","");let t=e.length===3,r=e.length===6,n=e.length===8;if(!t&&!r&&!n)throw new Error("unexpected hex "+e);let o=0,a=0,s=0;return t?(o=et(e.slice(0,1).repeat(2)),a=et(e.slice(1,2).repeat(2)),s=et(e.slice(2,3).repeat(2))):r?(o=et(e.slice(0,2)),a=et(e.slice(2,4)),s=et(e.slice(4,6))):n&&(o=et(e.slice(2,4)),a=et(e.slice(4,6)),s=et(e.slice(6,8))),(255<<24|(o&255)<<16|(a&255)<<8|s&255)>>>0}function et(e){return parseInt(e,16)}function le(e){return Ot(this,null,function*(){let t=yield new Promise((s,i)=>{let l=document.createElement("canvas"),u=l.getContext("2d");if(!u){i(new Error("Could not get canvas context"));return}let h=()=>{l.width=e.width,l.height=e.height,u.drawImage(e,0,0);let p=[0,0,e.width,e.height],y=e.dataset.area;y&&/^\d+(\s*,\s*\d+){3}$/.test(y)&&(p=y.split(/\s*,\s*/).map(k=>parseInt(k,10)));let[m,d,f,b]=p;s(u.getImageData(m,d,f,b).data)};e.complete?h():e.onload=h}),r=[];for(let s=0;s<t.length;s+=4){let i=t[s],l=t[s+1],u=t[s+2];if(t[s+3]<255)continue;let p=pt(i,l,u);r.push(p)}let n=Ft.quantize(r,128);return Q.score(n)[0]})}function Ce(e,t=[]){let r=G.of(e);return{source:e,schemes:{light:Pt.light(e),dark:Pt.dark(e)},palettes:{primary:r.a1,secondary:r.a2,tertiary:r.a3,neutral:r.n1,neutralVariant:r.n2,error:r.error},customColors:t.map(n=>be(e,n))}}function Ne(r){return Ot(this,arguments,function*(e,t=[]){let n=yield le(e);return Ce(n,t)})}function be(e,t){let r=t.value,n=r,o=e;t.blend&&(r=It.harmonize(n,o));let s=G.of(r).a1;return{color:t,value:r,light:{color:s.tone(40),onColor:s.tone(100),colorContainer:s.tone(90),onColorContainer:s.tone(10)},dark:{color:s.tone(80),onColor:s.tone(20),colorContainer:s.tone(30),onColorContainer:s.tone(90)}}}function Ue(e,t){var a,s;let r=(t==null?void 0:t.target)||document.body,o=((a=t==null?void 0:t.dark)!=null?a:!1)?e.schemes.dark:e.schemes.light;if(he(r,o),t!=null&&t.brightnessSuffix&&(he(r,e.schemes.dark,"-dark"),he(r,e.schemes.light,"-light")),t!=null&&t.paletteTones){let i=(s=t==null?void 0:t.paletteTones)!=null?s:[];for(let[l,u]of Object.entries(e.palettes)){let h=l.replace(/([a-z])([A-Z])/g,"$1-$2").toLowerCase();for(let p of i){let y=`--md-ref-palette-${h}-${h}${p}`,m=Ut(u.tone(p));r.style.setProperty(y,m)}}}}function he(e,t,r=""){for(let[n,o]of Object.entries(t.toJSON())){let a=n.replace(/([a-z])([A-Z])/g,"$1-$2").toLowerCase(),s=Ut(o);e.style.setProperty(`--md-sys-color-${a}${r}`,s)}}Object.assign(globalThis,ue);})();
/*! Bundled license information:

@material/material-color-utilities/utils/math_utils.js:
  (**
   * @license
   * Copyright 2021 Google LLC
   *
   * Licensed under the Apache License, Version 2.0 (the "License");
   * you may not use this file except in compliance with the License.
   * You may obtain a copy of the License at
   *
   *      http://www.apache.org/licenses/LICENSE-2.0
   *
   * Unless required by applicable law or agreed to in writing, software
   * distributed under the License is distributed on an "AS IS" BASIS,
   * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   * See the License for the specific language governing permissions and
   * limitations under the License.
   *)

@material/material-color-utilities/utils/color_utils.js:
  (**
   * @license
   * Copyright 2021 Google LLC
   *
   * Licensed under the Apache License, Version 2.0 (the "License");
   * you may not use this file except in compliance with the License.
   * You may obtain a copy of the License at
   *
   *      http://www.apache.org/licenses/LICENSE-2.0
   *
   * Unless required by applicable law or agreed to in writing, software
   * distributed under the License is distributed on an "AS IS" BASIS,
   * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   * See the License for the specific language governing permissions and
   * limitations under the License.
   *)

@material/material-color-utilities/hct/viewing_conditions.js:
  (**
   * @license
   * Copyright 2021 Google LLC
   *
   * Licensed under the Apache License, Version 2.0 (the "License");
   * you may not use this file except in compliance with the License.
   * You may obtain a copy of the License at
   *
   *      http://www.apache.org/licenses/LICENSE-2.0
   *
   * Unless required by applicable law or agreed to in writing, software
   * distributed under the License is distributed on an "AS IS" BASIS,
   * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   * See the License for the specific language governing permissions and
   * limitations under the License.
   *)

@material/material-color-utilities/hct/cam16.js:
  (**
   * @license
   * Copyright 2021 Google LLC
   *
   * Licensed under the Apache License, Version 2.0 (the "License");
   * you may not use this file except in compliance with the License.
   * You may obtain a copy of the License at
   *
   *      http://www.apache.org/licenses/LICENSE-2.0
   *
   * Unless required by applicable law or agreed to in writing, software
   * distributed under the License is distributed on an "AS IS" BASIS,
   * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   * See the License for the specific language governing permissions and
   * limitations under the License.
   *)

@material/material-color-utilities/hct/hct_solver.js:
  (**
   * @license
   * Copyright 2021 Google LLC
   *
   * Licensed under the Apache License, Version 2.0 (the "License");
   * you may not use this file except in compliance with the License.
   * You may obtain a copy of the License at
   *
   *      http://www.apache.org/licenses/LICENSE-2.0
   *
   * Unless required by applicable law or agreed to in writing, software
   * distributed under the License is distributed on an "AS IS" BASIS,
   * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   * See the License for the specific language governing permissions and
   * limitations under the License.
   *)

@material/material-color-utilities/hct/hct.js:
  (**
   * @license
   * Copyright 2021 Google LLC
   *
   * Licensed under the Apache License, Version 2.0 (the "License");
   * you may not use this file except in compliance with the License.
   * You may obtain a copy of the License at
   *
   *      http://www.apache.org/licenses/LICENSE-2.0
   *
   * Unless required by applicable law or agreed to in writing, software
   * distributed under the License is distributed on an "AS IS" BASIS,
   * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   * See the License for the specific language governing permissions and
   * limitations under the License.
   *)

@material/material-color-utilities/blend/blend.js:
  (**
   * @license
   * Copyright 2021 Google LLC
   *
   * Licensed under the Apache License, Version 2.0 (the "License");
   * you may not use this file except in compliance with the License.
   * You may obtain a copy of the License at
   *
   *      http://www.apache.org/licenses/LICENSE-2.0
   *
   * Unless required by applicable law or agreed to in writing, software
   * distributed under the License is distributed on an "AS IS" BASIS,
   * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   * See the License for the specific language governing permissions and
   * limitations under the License.
   *)

@material/material-color-utilities/contrast/contrast.js:
  (**
   * @license
   * Copyright 2022 Google LLC
   *
   * Licensed under the Apache License, Version 2.0 (the "License");
   * you may not use this file except in compliance with the License.
   * You may obtain a copy of the License at
   *
   *      http://www.apache.org/licenses/LICENSE-2.0
   *
   * Unless required by applicable law or agreed to in writing, software
   * distributed under the License is distributed on an "AS IS" BASIS,
   * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   * See the License for the specific language governing permissions and
   * limitations under the License.
   *)

@material/material-color-utilities/dislike/dislike_analyzer.js:
  (**
   * @license
   * Copyright 2023 Google LLC
   *
   * Licensed under the Apache License, Version 2.0 (the "License");
   * you may not use this file except in compliance with the License.
   * You may obtain a copy of the License at
   *
   *      http://www.apache.org/licenses/LICENSE-2.0
   *
   * Unless required by applicable law or agreed to in writing, software
   * distributed under the License is distributed on an "AS IS" BASIS,
   * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   * See the License for the specific language governing permissions and
   * limitations under the License.
   *)

@material/material-color-utilities/dynamiccolor/dynamic_color.js:
  (**
   * @license
   * Copyright 2022 Google LLC
   *
   * Licensed under the Apache License, Version 2.0 (the "License");
   * you may not use this file except in compliance with the License.
   * You may obtain a copy of the License at
   *
   *      http://www.apache.org/licenses/LICENSE-2.0
   *
   * Unless required by applicable law or agreed to in writing, software
   * distributed under the License is distributed on an "AS IS" BASIS,
   * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   * See the License for the specific language governing permissions and
   * limitations under the License.
   *)

@material/material-color-utilities/scheme/variant.js:
  (**
   * @license
   * Copyright 2022 Google LLC
   *
   * Licensed under the Apache License, Version 2.0 (the "License");
   * you may not use this file except in compliance with the License.
   * You may obtain a copy of the License at
   *
   *      http://www.apache.org/licenses/LICENSE-2.0
   *
   * Unless required by applicable law or agreed to in writing, software
   * distributed under the License is distributed on an "AS IS" BASIS,
   * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   * See the License for the specific language governing permissions and
   * limitations under the License.
   *)

@material/material-color-utilities/dynamiccolor/contrast_curve.js:
  (**
   * @license
   * Copyright 2023 Google LLC
   *
   * Licensed under the Apache License, Version 2.0 (the "License");
   * you may not use this file except in compliance with the License.
   * You may obtain a copy of the License at
   *
   *      http://www.apache.org/licenses/LICENSE-2.0
   *
   * Unless required by applicable law or agreed to in writing, software
   * distributed under the License is distributed on an "AS IS" BASIS,
   * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   * See the License for the specific language governing permissions and
   * limitations under the License.
   *)

@material/material-color-utilities/dynamiccolor/tone_delta_pair.js:
  (**
   * @license
   * Copyright 2023 Google LLC
   *
   * Licensed under the Apache License, Version 2.0 (the "License");
   * you may not use this file except in compliance with the License.
   * You may obtain a copy of the License at
   *
   *      http://www.apache.org/licenses/LICENSE-2.0
   *
   * Unless required by applicable law or agreed to in writing, software
   * distributed under the License is distributed on an "AS IS" BASIS,
   * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   * See the License for the specific language governing permissions and
   * limitations under the License.
   *)

@material/material-color-utilities/dynamiccolor/material_dynamic_colors.js:
  (**
   * @license
   * Copyright 2022 Google LLC
   *
   * Licensed under the Apache License, Version 2.0 (the "License");
   * you may not use this file except in compliance with the License.
   * You may obtain a copy of the License at
   *
   *      http://www.apache.org/licenses/LICENSE-2.0
   *
   * Unless required by applicable law or agreed to in writing, software
   * distributed under the License is distributed on an "AS IS" BASIS,
   * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   * See the License for the specific language governing permissions and
   * limitations under the License.
   *)

@material/material-color-utilities/palettes/tonal_palette.js:
  (**
   * @license
   * Copyright 2021 Google LLC
   *
   * Licensed under the Apache License, Version 2.0 (the "License");
   * you may not use this file except in compliance with the License.
   * You may obtain a copy of the License at
   *
   *      http://www.apache.org/licenses/LICENSE-2.0
   *
   * Unless required by applicable law or agreed to in writing, software
   * distributed under the License is distributed on an "AS IS" BASIS,
   * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   * See the License for the specific language governing permissions and
   * limitations under the License.
   *)

@material/material-color-utilities/palettes/core_palette.js:
  (**
   * @license
   * Copyright 2021 Google LLC
   *
   * Licensed under the Apache License, Version 2.0 (the "License");
   * you may not use this file except in compliance with the License.
   * You may obtain a copy of the License at
   *
   *      http://www.apache.org/licenses/LICENSE-2.0
   *
   * Unless required by applicable law or agreed to in writing, software
   * distributed under the License is distributed on an "AS IS" BASIS,
   * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   * See the License for the specific language governing permissions and
   * limitations under the License.
   *)

@material/material-color-utilities/quantize/lab_point_provider.js:
  (**
   * @license
   * Copyright 2021 Google LLC
   *
   * Licensed under the Apache License, Version 2.0 (the "License");
   * you may not use this file except in compliance with the License.
   * You may obtain a copy of the License at
   *
   *      http://www.apache.org/licenses/LICENSE-2.0
   *
   * Unless required by applicable law or agreed to in writing, software
   * distributed under the License is distributed on an "AS IS" BASIS,
   * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   * See the License for the specific language governing permissions and
   * limitations under the License.
   *)

@material/material-color-utilities/quantize/quantizer_wsmeans.js:
  (**
   * @license
   * Copyright 2021 Google LLC
   *
   * Licensed under the Apache License, Version 2.0 (the "License");
   * you may not use this file except in compliance with the License.
   * You may obtain a copy of the License at
   *
   *      http://www.apache.org/licenses/LICENSE-2.0
   *
   * Unless required by applicable law or agreed to in writing, software
   * distributed under the License is distributed on an "AS IS" BASIS,
   * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   * See the License for the specific language governing permissions and
   * limitations under the License.
   *)

@material/material-color-utilities/quantize/quantizer_map.js:
  (**
   * @license
   * Copyright 2021 Google LLC
   *
   * Licensed under the Apache License, Version 2.0 (the "License");
   * you may not use this file except in compliance with the License.
   * You may obtain a copy of the License at
   *
   *      http://www.apache.org/licenses/LICENSE-2.0
   *
   * Unless required by applicable law or agreed to in writing, software
   * distributed under the License is distributed on an "AS IS" BASIS,
   * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   * See the License for the specific language governing permissions and
   * limitations under the License.
   *)

@material/material-color-utilities/quantize/quantizer_wu.js:
  (**
   * @license
   * Copyright 2021 Google LLC
   *
   * Licensed under the Apache License, Version 2.0 (the "License");
   * you may not use this file except in compliance with the License.
   * You may obtain a copy of the License at
   *
   *      http://www.apache.org/licenses/LICENSE-2.0
   *
   * Unless required by applicable law or agreed to in writing, software
   * distributed under the License is distributed on an "AS IS" BASIS,
   * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   * See the License for the specific language governing permissions and
   * limitations under the License.
   *)

@material/material-color-utilities/quantize/quantizer_celebi.js:
  (**
   * @license
   * Copyright 2021 Google LLC
   *
   * Licensed under the Apache License, Version 2.0 (the "License");
   * you may not use this file except in compliance with the License.
   * You may obtain a copy of the License at
   *
   *      http://www.apache.org/licenses/LICENSE-2.0
   *
   * Unless required by applicable law or agreed to in writing, software
   * distributed under the License is distributed on an "AS IS" BASIS,
   * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   * See the License for the specific language governing permissions and
   * limitations under the License.
   *)

@material/material-color-utilities/scheme/dynamic_scheme.js:
  (**
   * @license
   * Copyright 2022 Google LLC
   *
   * Licensed under the Apache License, Version 2.0 (the "License");
   * you may not use this file except in compliance with the License.
   * You may obtain a copy of the License at
   *
   *      http://www.apache.org/licenses/LICENSE-2.0
   *
   * Unless required by applicable law or agreed to in writing, software
   * distributed under the License is distributed on an "AS IS" BASIS,
   * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   * See the License for the specific language governing permissions and
   * limitations under the License.
   *)

@material/material-color-utilities/scheme/scheme.js:
  (**
   * @license
   * Copyright 2021 Google LLC
   *
   * Licensed under the Apache License, Version 2.0 (the "License");
   * you may not use this file except in compliance with the License.
   * You may obtain a copy of the License at
   *
   *      http://www.apache.org/licenses/LICENSE-2.0
   *
   * Unless required by applicable law or agreed to in writing, software
   * distributed under the License is distributed on an "AS IS" BASIS,
   * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   * See the License for the specific language governing permissions and
   * limitations under the License.
   *)

@material/material-color-utilities/scheme/scheme_android.js:
  (**
   * @license
   * Copyright 2021 Google LLC
   *
   * Licensed under the Apache License, Version 2.0 (the "License");
   * you may not use this file except in compliance with the License.
   * You may obtain a copy of the License at
   *
   *      http://www.apache.org/licenses/LICENSE-2.0
   *
   * Unless required by applicable law or agreed to in writing, software
   * distributed under the License is distributed on an "AS IS" BASIS,
   * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   * See the License for the specific language governing permissions and
   * limitations under the License.
   *)

@material/material-color-utilities/temperature/temperature_cache.js:
  (**
   * @license
   * Copyright 2023 Google LLC
   *
   * Licensed under the Apache License, Version 2.0 (the "License");
   * you may not use this file except in compliance with the License.
   * You may obtain a copy of the License at
   *
   *      http://www.apache.org/licenses/LICENSE-2.0
   *
   * Unless required by applicable law or agreed to in writing, software
   * distributed under the License is distributed on an "AS IS" BASIS,
   * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   * See the License for the specific language governing permissions and
   * limitations under the License.
   *)

@material/material-color-utilities/scheme/scheme_content.js:
  (**
   * @license
   * Copyright 2023 Google LLC
   *
   * Licensed under the Apache License, Version 2.0 (the "License");
   * you may not use this file except in compliance with the License.
   * You may obtain a copy of the License at
   *
   *      http://www.apache.org/licenses/LICENSE-2.0
   *
   * Unless required by applicable law or agreed to in writing, software
   * distributed under the License is distributed on an "AS IS" BASIS,
   * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   * See the License for the specific language governing permissions and
   * limitations under the License.
   *)

@material/material-color-utilities/scheme/scheme_expressive.js:
  (**
   * @license
   * Copyright 2022 Google LLC
   *
   * Licensed under the Apache License, Version 2.0 (the "License");
   * you may not use this file except in compliance with the License.
   * You may obtain a copy of the License at
   *
   *      http://www.apache.org/licenses/LICENSE-2.0
   *
   * Unless required by applicable law or agreed to in writing, software
   * distributed under the License is distributed on an "AS IS" BASIS,
   * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   * See the License for the specific language governing permissions and
   * limitations under the License.
   *)

@material/material-color-utilities/scheme/scheme_fidelity.js:
  (**
   * @license
   * Copyright 2023 Google LLC
   *
   * Licensed under the Apache License, Version 2.0 (the "License");
   * you may not use this file except in compliance with the License.
   * You may obtain a copy of the License at
   *
   *      http://www.apache.org/licenses/LICENSE-2.0
   *
   * Unless required by applicable law or agreed to in writing, software
   * distributed under the License is distributed on an "AS IS" BASIS,
   * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   * See the License for the specific language governing permissions and
   * limitations under the License.
   *)

@material/material-color-utilities/scheme/scheme_monochrome.js:
  (**
   * @license
   * Copyright 2022 Google LLC
   *
   * Licensed under the Apache License, Version 2.0 (the "License");
   * you may not use this file except in compliance with the License.
   * You may obtain a copy of the License at
   *
   *      http://www.apache.org/licenses/LICENSE-2.0
   *
   * Unless required by applicable law or agreed to in writing, software
   * distributed under the License is distributed on an "AS IS" BASIS,
   * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   * See the License for the specific language governing permissions and
   * limitations under the License.
   *)

@material/material-color-utilities/scheme/scheme_neutral.js:
  (**
   * @license
   * Copyright 2022 Google LLC
   *
   * Licensed under the Apache License, Version 2.0 (the "License");
   * you may not use this file except in compliance with the License.
   * You may obtain a copy of the License at
   *
   *      http://www.apache.org/licenses/LICENSE-2.0
   *
   * Unless required by applicable law or agreed to in writing, software
   * distributed under the License is distributed on an "AS IS" BASIS,
   * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   * See the License for the specific language governing permissions and
   * limitations under the License.
   *)

@material/material-color-utilities/scheme/scheme_tonal_spot.js:
  (**
   * @license
   * Copyright 2022 Google LLC
   *
   * Licensed under the Apache License, Version 2.0 (the "License");
   * you may not use this file except in compliance with the License.
   * You may obtain a copy of the License at
   *
   *      http://www.apache.org/licenses/LICENSE-2.0
   *
   * Unless required by applicable law or agreed to in writing, software
   * distributed under the License is distributed on an "AS IS" BASIS,
   * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   * See the License for the specific language governing permissions and
   * limitations under the License.
   *)

@material/material-color-utilities/scheme/scheme_vibrant.js:
  (**
   * @license
   * Copyright 2022 Google LLC
   *
   * Licensed under the Apache License, Version 2.0 (the "License");
   * you may not use this file except in compliance with the License.
   * You may obtain a copy of the License at
   *
   *      http://www.apache.org/licenses/LICENSE-2.0
   *
   * Unless required by applicable law or agreed to in writing, software
   * distributed under the License is distributed on an "AS IS" BASIS,
   * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   * See the License for the specific language governing permissions and
   * limitations under the License.
   *)

@material/material-color-utilities/score/score.js:
  (**
   * @license
   * Copyright 2021 Google LLC
   *
   * Licensed under the Apache License, Version 2.0 (the "License");
   * you may not use this file except in compliance with the License.
   * You may obtain a copy of the License at
   *
   *      http://www.apache.org/licenses/LICENSE-2.0
   *
   * Unless required by applicable law or agreed to in writing, software
   * distributed under the License is distributed on an "AS IS" BASIS,
   * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   * See the License for the specific language governing permissions and
   * limitations under the License.
   *)

@material/material-color-utilities/utils/string_utils.js:
  (**
   * @license
   * Copyright 2021 Google LLC
   *
   * Licensed under the Apache License, Version 2.0 (the "License");
   * you may not use this file except in compliance with the License.
   * You may obtain a copy of the License at
   *
   *      http://www.apache.org/licenses/LICENSE-2.0
   *
   * Unless required by applicable law or agreed to in writing, software
   * distributed under the License is distributed on an "AS IS" BASIS,
   * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   * See the License for the specific language governing permissions and
   * limitations under the License.
   *)

@material/material-color-utilities/utils/image_utils.js:
  (**
   * @license
   * Copyright 2021 Google LLC
   *
   * Licensed under the Apache License, Version 2.0 (the "License");
   * you may not use this file except in compliance with the License.
   * You may obtain a copy of the License at
   *
   *      http://www.apache.org/licenses/LICENSE-2.0
   *
   * Unless required by applicable law or agreed to in writing, software
   * distributed under the License is distributed on an "AS IS" BASIS,
   * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   * See the License for the specific language governing permissions and
   * limitations under the License.
   *)

@material/material-color-utilities/utils/theme_utils.js:
  (**
   * @license
   * Copyright 2021 Google LLC
   *
   * Licensed under the Apache License, Version 2.0 (the "License");
   * you may not use this file except in compliance with the License.
   * You may obtain a copy of the License at
   *
   *      http://www.apache.org/licenses/LICENSE-2.0
   *
   * Unless required by applicable law or agreed to in writing, software
   * distributed under the License is distributed on an "AS IS" BASIS,
   * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   * See the License for the specific language governing permissions and
   * limitations under the License.
   *)

@material/material-color-utilities/index.js:
  (**
   * @license
   * Copyright 2021 Google LLC
   *
   * Licensed under the Apache License, Version 2.0 (the "License");
   * you may not use this file except in compliance with the License.
   * You may obtain a copy of the License at
   *
   *      http://www.apache.org/licenses/LICENSE-2.0
   *
   * Unless required by applicable law or agreed to in writing, software
   * distributed under the License is distributed on an "AS IS" BASIS,
   * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   * See the License for the specific language governing permissions and
   * limitations under the License.
   *)
*/
