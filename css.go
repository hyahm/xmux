package xmux

var style = `
*{box-sizing: border-box}
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
    padding: 10px 0;
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
