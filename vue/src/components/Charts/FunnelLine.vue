<template>
  <div id="chartDiv" />
</template>
<script>
import * as G2 from '@antv/g2' // 引入G2
import { Chart } from '@antv/g2' // 或者只引入需要用到的G2组件，我要用Chart
const DataSet = require('@antv/data-set') // antV中用DataSet作为数据集，可以对原始数据进行操作，比如格式转换之类的，以供图表方法使用，当然也可以不用，不管什么方法只要把数据格式转变成它需要的格式就可以了
import failPng from '@/assets/img/fail.png'
import toPng from '@/assets/img/to.png'

export default {
  name: 'FunnelLine',
  props: {
    chartData: {
      type: Array,
      default: []
    }
  },
  data() {
    return {

      failPng: failPng,
      toPng: toPng
    }
  },
  destroyed() {

  },
  mounted() {
    this.initLineChart()
  },

  methods: {
    // 图表
    initLineChart() {
      if (this.chartData.length == 0) {
        return
      }
      const chart = new Chart({ // 创建一个图表
        container: 'chartDiv', // 容器是上面那个div
        autoFit: true, // 自适应
        height: 500, // 高度
        padding: [60, 20, 40, 60]
      })

      var colorlists = [
        'hsl(42, 48%, 54%)',
        'hsl(138, 24%, 48%)',
        'rgb(200, 138, 131)',
        'rgb(84, 221, 226)',
        'rgb(178, 199, 168)',
        'rgb(16, 195, 195)',
        'hsl(0, 21%, 68%)',
        'rgb(226, 166, 198)',
        'hsl(278, 17%, 66%)',
        'rgb(153, 199, 235)',
        'blueviolet']

      var Shape = G2
      Shape.registerShape('interval', 'textInterval', {
        draw: function(cfg, group) {
          var points = this.parsePoints(cfg.points) // 将0-1空间的坐标转换为画布坐标
          var count = cfg.data.count
          var succScale = cfg.data.succScale
          group.addShape('text', {
            attrs: {
              text: count,
              textAlign: 'center',
              x: points[1].x + cfg.size / 2,
              y: points[1].y,
              fontFamily: 'PingFang SC',
              fontSize: 12,
              fill: '#BBB'
            }
          })

          var polygon = group.addShape('polygon', {
            attrs: {
              points: points.map((point) => {
                return [point.x, point.y]
              }),
              fill: cfg.color
            }
          })
          return polygon
        }
      })
      const tmp = this.chartData[this.chartData.length - 1]

      Shape.registerShape('interval', 'fallFlag', {
        getPoints: (_ref) => {
          var x = _ref.x
          var y = _ref.y
          var y0 = _ref.y0
          var size = _ref.size

          return [{
            x: x + size,
            y: y0 + size
          }, {
            x: x,
            y: y
          }]
        },
        draw: function(cfg, group) {
          if (cfg.data === tmp) {
            return
          }
          var points = this.parsePoints(cfg.points) // 将0-1空间的坐标转换为画布坐标
          var p1 = points[0]
          var width = 9
          var lostScale = cfg.data.lostScale
          var succScale = cfg.data.succScale
          group.addShape('text', {
            attrs: {
              text: lostScale + ' %',
              x: p1.x - width / 2 - 14,
              y: p1.y - 14,
              fontFamily: 'PingFang SC',
              fontSize: 12,
              fill: '#BBB'
            }
          })
          group.addShape('text', {

            attrs: {
              text: succScale + '%',
              x: p1.x - width / 2 - 14,
              y: p1.y - 270,
              fontFamily: 'PingFang SC',
              fontSize: 12,
              fill: '#BBB'
            }
          })
          group.addShape('image', {
            attrs: {
              x: p1.x - 16,
              y: p1.y - 272,
              img: toPng,
              width: 40,
              height: 32
            }
          })
          var polygon = group.addShape('image', {
            attrs: {
              x: p1.x - 16,
              y: p1.y - 16,
              img: failPng,
              width: 40,
              height: 32
            }
          })
          return polygon // 将自定义Shape返回
        }
      })

      chart.legend(false)

      chart.source(this.chartData, {
        count: {
          alias: '访问数'
        },
        showTitle: {
          alias: '步骤名称'
        }
      })
      chart.axis('showTitle', {
        title: null
      })
      chart.interval().position('showTitle*max').color('#D5D5D5').size(30)　　// 第一次绘图，绘制背景，color设置颜色，opacity设置透明度
      chart.interval().position('showTitle*count').shape('textInterval').color('showTitle', (count) => {
        return colorlists[Math.floor(Math.random() * colorlists.length)]
      }).size(30)
      chart.interval().position('showTitle*count').color('#E4E4E4').shape('fallFlag')

      chart.render()
    }
  }
}
</script>
