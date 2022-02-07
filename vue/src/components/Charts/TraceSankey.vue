<template>
  <div :class="className" :id="className" :style="{height:height,width:width}"/>
</template>

<script>
  // import echarts from 'echarts'
  //require('echarts/theme/macarons') // echarts theme
  import resize from './mixins/resize'
  import {setNotopt} from './mixins/utils'

  export default {
    mixins: [resize],
    props: {
      className: {
        type: String,
        default: 'chart'
      },
      width: {
        type: String,
        default: '100%'
      },
      height: {
        type: String,
        default: '400px'
      },
      autoResize: {
        type: Boolean,
        default: true
      },
      chartData: {
        type: Array,
        default: []
      }
    },
    data() {
      return {
        chart: null
      }
    },
    watch: {},
    mounted() {
      this.initChart()
    },
    beforeDestroy() {
      if (!this.chart) {
        return
      }
      this.chart.dispose()
      this.chart = null
    },
    methods: {
      initChart() {
        this.chart = echarts.init(document.getElementById(this.className))

        this.setOptions()

      },
      getColor1() {//固定红色值
        var re = "#";
        var col = this.color();
        re += col + "FF";
        return re
      },
      getColor2() {//固定蓝色值
        var re = "#FF";
        var col = this.color();
        re += col;
        return re
      },

      color() {
        var re = "";
        for (var loopNum = 0; loopNum < 2; loopNum++) {
          var temp = Math.floor(256 * Math.random());
          if (temp < 130 && loopNum == 0) {
            temp = 130;
          }
          if (temp > 200 && loopNum == 1) {
            temp = 200;
          }
          temp = temp.toString(16);//将数值转换成16进制
          if (temp.length !== 2) {
            temp = "0" + temp
          }
          re += temp//对颜色进行拼接
        }
        return re;
      },

      setOptions() {
        if (this.chartData.length == 0) {
          setNotopt(this.chart)
          return
        }

        let eventArr = []
        let targetArr = []
        let eventSet = new Map()
        for (let k in this.chartData) {
          let traceCharts = this.chartData[k]

          eventSet.set(traceCharts['event'][0], 1)
          eventSet.set(traceCharts['event'][1], 1)

          targetArr.push({
            source: traceCharts['event'][0],
            target: traceCharts['event'][1],
            value: traceCharts['sum_user_count'],
          })
        }
        eventSet.forEach((v, k, tmp) => {
          let color = ""
          var random = Math.random();
          if(random <0.618){//分配红色和蓝色出现的比例
            color = this.getColor1()
          }else{
            color = this.getColor2()
          }
          let obj = {
            name: k,
            itemStyle: {
              color: color
            }
          }
          eventArr.push(obj)
        })

        let levelsArr = []

        for (let i in eventArr) {
          levelsArr.push({
            depth: i,
            itemStyle: {
              color: eventArr[i].itemStyle.color
            },
            lineStyle: {
              color: 'source',
              opacity: 0.4
            }
          })
        }

        let option = {
          tooltip: {
            trigger: 'item',
            triggerOn: 'mousemove'
          },
          backgroundColor: '#FFFFFF',
          series: {
            type: 'sankey',
            layout: 'none',
            top: 50,
            left: '3%',
            right: '12%',
            nodeGap: 14,
            layoutIterations: 0, // 自动优化列表，尽量减少线的交叉，为0就是按照数据排列
            data: eventArr, // 节点
            links: targetArr, // 节点之间的连线

            focusNodeAdjacency: 'allEdges', // 鼠标划上时高亮的节点和连线，allEdges表示鼠标划到节点上点亮节点上的连线及连线对应的节点
            levels: [{
              depth: 0,
              itemStyle: {
                color: '#F27E7E'
              },
              lineStyle: {
                color: 'source',
                opacity: 0.4
              }
            },
              {
                depth: 1,
                lineStyle: {
                  color: 'source',
                  opacity: 0.4
                }
              },
              {
                depth: 2,
                lineStyle: {
                  color: 'source',
                  opacity: 0.4
                }
              },
              {
                depth: 3,
                label: {
                  fontSize: 12
                }
              }
            ],
            label: {
              fontSize: 14,
              color: '#666'
            },
            itemStyle: {
              normal: {
                borderWidth: 0
              }
            }
          }
        }

          this.chart.setOption(
            option
          )


      }

    }
  }
</script>
