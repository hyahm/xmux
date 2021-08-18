package xmux

var style = `

.body-content{
    width: 1300px;
    margin: 0 auto;
    position:relative;
    background: #eee;
    height: 100%;
}
.body-head{
    width: 100%;
    color: #000;
    font-size: 20px;
    font-weight: bold;
    height: 60px;
    line-height: 60px;
    text-indent: 30px;
    position: fixed;
    left: calc((100% - 1300px) / 2);
    right: calc((100% - 1300px) / 2);
    z-index: 9999;
    background: #fff;
}
.left-search{
    padding:0 10px;
    background: #fff;
    display: flex;
    justify-content: space-around;
    align-items: center;
    border-radius: 20px;
    width: 90%;
    margin:  0 auto;
}

.input{
    width: 100%;
    height: 30px;
    background: #fff;
    border: 0;
    text-indent: 20px;
}
/*淇敼婊氬姩鏉℃牱寮�*/
.body-left::-webkit-scrollbar{
    width:10px;
    height:10px;
    /**/
}
.body-left::-webkit-scrollbar-track{
    background: rgb(239, 239, 239);
    border-radius:2px;
}
.body-left::-webkit-scrollbar-thumb{
    background: #bfbfbf;
    border-radius:10px;
}
.body-left::-webkit-scrollbar-thumb:hover{
    background: #333;
}
.body-left::-webkit-scrollbar-corner{
    background: #179a16;
}
.body-right{
    padding-top: 80px;
    width: calc(1300px - 240px);
    float: right;
}
.text-center{
    text-align: center;
}
.right-msg{
    padding-left: 30px;
}
.right-msg h4{
    padding-left: 20px;
}
.right-light{
    padding: 15px 30px 0;
}
.right-dl{
    background: #b2b8be;
    border-radius: 5px;
    height: 30px;
    padding: 5px 10px;
    display: flex;
    justify-content: flex-start;
}
.right-get{
    background:#c55b03;
    color: #fff;
    border-radius: 5px;
    display: block;
    width: 60px;
    text-align: center;
}
.right-post{
    background:#008cff;
}
.right-url{
    padding-left: 20px;
}
.dl-box{
    background: #fff;
    padding: 0 0 50px 0;
}
.dl-table{
    width: 100%;
    display: flex;
    justify-content: space-around;
    padding: 0 30px;
}
.dl-box h3{
     padding: 10px 20px;
}
.dl-table span{
    width:calc(100% / 5);
    text-align: center;
    background:#c55b03;
    padding: 10px 0;
    color: #fff;
    border-right: 1px solid #eee;
}
.dl-table span:first-child{
    border-left: 1px solid #eee;
}
.dl-table-msg span{
    width:calc(100% / 5);
    text-align: center;
    background:#fff;
    padding: 10px 0;
    color: #000;
    border-right: 1px solid #eee;
    border-bottom: 1px solid #eee;
}
.dl-ex-box{
    width: 100%;
    padding: 10px 40px;
}
.dl-expl{
    background: #e1e1e8;
    border: 1px solid #eee;
    text-align: left;
    padding: 10px;
    white-space: pre-wrap;
    color: #080;
}
.dl-table1 span{
    width:calc(100% / 3);
    text-align: center;
    background:#c55b03;
    padding: 10px 0;
    color: #fff;
    border-right: 1px solid #eee;
}
.dl-table-msg1 span{
    width:calc(100% / 3);
    text-align: center;
    background:#fff;
    padding: 10px 0;
    color: #000;
    border-right: 1px solid #eee;
    border-bottom: 1px solid #eee;
}
.dl-bz{
    padding: 0 40px 0;
}
.dl-none{
    display: none;
}
.dl-block{
    display: block;
}
.dl-post span{
    background:#008cff;
}`

var cssleft = `
@charset "utf-8";
/* 以下实际使用若已初始化可删除 .lsm-sidebar height父级需逐级设置为100%*/
*{
    margin: 0;
    padding: 0;
    box-sizing: border-box;
}
body,html{height: 100%;}
body,ul{margin:0;padding:0}
body,pre{font:14px "微软雅黑","宋体","Arial Narrow",HELVETICA;-webkit-text-size-adjust:100%;}
li{list-style:none}
a{text-decoration:none;}
input {

    border: 0;
outline: none;
background-color: rgba(0, 0, 0, 0);

}

.left-side-menu,.lsm-popup.lsm-sidebar   ul li, .lsm-container.lsm-mini .lsm-sidebar>ul>li.lsm-sidebar-item>ul>li.lsm-sidebar-item>ul{
    background: #3b3e47;
}


.left-side-menu {-webkit-touch-callout: none;-webkit-user-select: none;-khtml-user-select: none;-moz-user-select: none;-ms-user-select: none; user-select: none; }
.left-side-menu{height: calc(100% - 60px);width: 240px; top: 60px;position: absolute; padding-top: 20px}
.lsm-expand-btn{height: 65px;}
.lsm-container {height: calc(100% - 65px); transition: all .3s;z-index: 100}

.lsm-container li>a.active{ background: #c55b03; color: #fff; }


.lsm-sidebar a{display: block;overflow: hidden;padding-left: 20px;line-height: 40px;max-height: 40px;color: #b2b8be;transition: all .3s;}
.lsm-container ul:first-child > li> a, .lsm-container ul:first-child > li> a span{ line-height: 55px;max-height: 55px; }
.lsm-sidebar a span{margin-left: 30px;}
.lsm-sidebar .lsm-sidebar-item .lsm-sidebar-item >ul>li a span{margin-left: 60px;}
.lsm-sidebar-item{position: relative;}
.lsm-sidebar-item.lsm-sidebar-show{border-bottom: none;}
.lsm-sidebar-item ul{display: none;background: rgba(0,0,0,.1);}
.lsm-sidebar-item.lsm-sidebar-show ul{display: block;}
.lsm-sidebar-item>a:before{content: "";position: absolute;left: 0px;width: 2px;height: 40px;background: #34A0CE;opacity:0;transition: all .3s;}
.lsm-container ul:first-child>li.lsm-sidebar-item>a:before{height: 55px;}
.lsm-sidebar .lsm-sidebar-icon{font-size: 20px;position: absolute;margin-left:-1px;}
/* 此处修改导航图标 可自定义iconfont 替换*/
.icon_1::after{content: "\e506";}

.lsm-sidebar-more{float:right;margin-right: 20px;font-size: 12px;transition: transform .3s;}

/* 导航右侧箭头 换用其他字体需要替换*/
.lsm-sidebar-more::after{content: "\e93d";}


.lsm-sidebar-show > a > i.my-icon.lsm-sidebar-more{transform:rotate(90deg);}
.lsm-sidebar-show,.lsm-sidebar-item>a:hover{color: #FFF;background: rgba(0, 0, 0, 0.2);}
.lsm-sidebar-show>a:before,.lsm-sidebar-item>a:hover:before{opacity:1;}
.lsm-sidebar-item li>a:hover,.lsm-popup>div>ul>li>a:hover{color: #FFF; background: #6e809c;}
.lsm-mini-btn{height: 70px;width: 70px;}
.lsm-mini-btn svg{margin: -10px 0 0 -10px;}
.lsm-mini-btn input[type="checkbox"]{display: none;}

.lsm-mini-btn path {
    fill: none;
    stroke: #ffffff;
    stroke-width: 3;
    stroke-linecap: round;
    stroke-linejoin: round;
    --length: 24;
    --offset: -38;
    stroke-dasharray: var(--length) var(--total-length);
    stroke-dashoffset: var(--offset);
    transition: all .8s cubic-bezier(.645, .045, .355, 1);
}

.lsm-mini-btn circle {fill: #fff3;opacity: 0;}
.lsm-mini-btn label {top: 0; right: 0;}
.lsm-mini-btn label:hover circle {opacity: 1;}
.lsm-mini-btn input:checked+svg .line--1, .lsm-mini-btn input:checked+svg .line--3 {--length: 8.602325267;}
.lsm-mini-btn .line--1, .lsm-mini-btn .line--3 {--total-length: 126.38166809082031;}
.lsm-mini-btn .line--2 {--total-length: 80;}
.lsm-mini-btn input:checked+svg .line--1, .lsm-mini-btn input:checked+svg .line--3 {--offset: -109.1770175568;}

.lsm-mini .lsm-container, .lsm-mini .lsm-container{width: 60px;}
.lsm-container.lsm-mini .lsm-sidebar .lsm-sidebar-icon{/* margin-left:-2px; */}
.left-side-menu.lsm-mini ul:first-child>li.lsm-sidebar-item>a span{display: none;}
.left-side-menu.lsm-mini ul:first-child>li.lsm-sidebar-item>a> i.lsm-sidebar-more{margin-right: -20px;}
.lsm-container.lsm-mini .lsm-sidebar>ul>li.lsm-sidebar-item>ul>li.lsm-sidebar-item>ul{
    display:none;
    position: absolute;top:0px;left:180px;width: 180px;z-index: 99;
    bottom: 0px;
    top: 0px;
    overflow: hidden;
}
.left-side-menu.lsm-mini ul:first-child > li > ul{
    display: none;
}
.transform { -webkit-transform: scale(1); -ms-transform: scale(1); transform: scale(1); }
.lsm-popup div{background: #05161f;}
.lsm-popup{
    display: block;
    position: absolute;
    border: 3px solid rgba(60, 71, 76, 0);
}

.lsm-popup > div > a > i.my-icon.lsm-sidebar-more{
    transform:rotate(90deg);
}

.lsm-popup.second{
    left: 60px;
}
.lsm-popup.third{
    left: 243px;
}
.lsm-popup.third.lsm-sidebar > div > ul {
    display: block;
}
.lsm-popup div {
    border-radius: 5px;
}
.lsm-popup .lsm-sidebar-icon{
    display: none;
}
.lsm-popup.lsm-sidebar a span{
    margin-left: 0px;
}
.lsm-popup.lsm-sidebar > div > ul > li.lsm-sidebar-item>ul{position: absolute;top:0px;left:180px;width: 180px;z-index: 99;}

.lsm-popup.lsm-sidebar   ul {
    width: 180px;
}
.lsm-popup.lsm-sidebar   ul li{
    width: 180px;
}
.lsm-popup.lsm-sidebar ul li:last-child, .lsm-popup>div>ul>li:last-child>a{
    border-radius: 0 0 5px 5px ;
}

`

var font = `
@font-face {font-family: "iconfont";
  src: url('iconfont.eot?t=1587720235809'); /* IE9 */
  src: url('iconfont.eot?t=1587720235809#iefix') format('embedded-opentype'), /* IE6-IE8 */
  url('data:application/x-font-woff2;charset=utf-8;base64,d09GMgABAAAAAAQsAAsAAAAACKAAAAPfAAEAAAAAAAAAAAAAAAAAAAAAAAAAAAAAHEIGVgCDMgqEOINpATYCJAMUCwwABCAFhG0HTht7B8ieGjsLiZpN0DQ0sKiZDVzO44F+v2/nypN1TaJxq9j09RBFkYpXfDpMZ0h0C5kmlvBQaP/IzmYT7z2Fookeglujd9Hmai43D1QEtcKUABXeJHsB3r8UXH2FiXzf8S9rgZxq0zpr7/vhfqaL0g+2evZtjlFpkVZtm01HCTSg6Lg7kjuQO1G+YezyIrcT0GycAbTptX8OJZneLyCcx5gCSgm1XBEa1IKyYGqG8BTy6lhD+gbwxPv58A9GAwNERQJ96NbhhgHLP2h/FsjF/GdXDgY4wxnBZpGwAGTC9ULzuWzgXUDWLP9dzirQrBbFH4WfDT/TryU+H4W/IWmz8Q+PJIsCoQQ9uBWsIlE+OI6FkvjZgJL5maJE/ixwr2WXAjSjDS+B74g40TSIKqmSa5Rv0pdBoppchZoWOR8RvhCVI9+e7j0CSu09pfhhOPHgaSWio66rTyr36Yn78MDOE7AqXpmH4dlrC7x5qITw+r4CZeDAUnjrSJk6eGjlL1DbWdJFNNV1FLdjTaAzhkgLdhRLoMJjQom6OXnUO6uFx6dloOcFoeVhVDKLGzmd1tNqpyF3d4Zoe5mS9uub1col2ILQfn9cgpXQ7CfMe6uob7Oob41OHPeGw+yqFlP1byizwiwn8OZptMqkBgWO5FRFBt0+jbjbx2iclTWXYm/aoRkJk/LGx3JthDQ2huSOjd9AGh+H0fhvC8LrzpvWfWg6EdlDCdoTccUJ0S5hw3/3f4ODfp40/DfwV8sbLOys+31dlA392OrQcl0AS6+ebpQKr4tEVxSLUtRbPIXHrVVEFBsVZtLckCO0qT19qCr1Z22lVj09f2jYy1Dfz9SAbICGSRpGNNIreaJv6QpZHYBckP+B33H6kfo/SNjhv26gAHx6v9GHygbNzJYAL1KiMeeNp2Wgcc5zPF1TZCnbYkudCvuKILKyZrxne7fPYeR999AY1AkOg6gyBkltApuJC1BosQil2go0m+c9u0WPVFkhM4E5FQCCTksQtXsHSacLbCZ+QKHfH5Q6AxmaHQTxJVtMhRjsnkpCkYGhORgTFgbS1KmKH3pKZuIJjvM86iVx1Xag6eBwOrBFAfEQI9IyzdH0dImRqaGPXmA38rwQk1NDh3Q1GJ2uolnjIVn0QYNhST5UbGBEUIgBhXqCYuhCAemUViqZ7z9FTAmPwBU5ZdZLhFXZjSNTRkwaYLfEQaOcS2md5p5pVDpF4l6RWMiHvKCReIY3hJKLOzmIThkU3cEdmWWMCkU2FQ9OL/XvcR0007+kCCRERhRERTTIF0lM5bCwgbQIwbpVzI4RgeknxUNWAQAA') format('woff2'),
  url('iconfont.woff?t=1587720235809') format('woff'),
  url('iconfont.ttf?t=1587720235809') format('truetype'), /* chrome, firefox, opera, Safari, Android, iOS 4.2+ */
  url('iconfont.svg?t=1587720235809#iconfont') format('svg'); /* iOS 4.1- */
}

.iconfont {
  font-family: "iconfont" !important;
  font-size: 16px;
  font-style: normal;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
}

.icon-more:before {
  content: "\e93d";
}

.icon-cc-search:before {
  content: "\e698";
}

.icon-xiangmu:before {
  content: "\e506";
}

.icon-project:before {
  content: "\e60b";
}`
