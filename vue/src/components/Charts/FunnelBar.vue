<template>
  <div :id="className" :class="className" :style="{height:height,width:width}" />
</template>

<script>
// import echarts from 'echarts'
// require('echarts/theme/macarons') // echarts theme
import resize from './mixins/resize'
import { setNotopt } from './mixins/utils'
import toPng from '@/assets/img/arrow-v.svg'

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
  watch: {

  },
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
    setOptions() {
      if (this.chartData.length == 0) {
        setNotopt(this.chart)
        return
      }
      let option

      const title = this.chartData.map(obj => { return obj.showTitle })
      const data = this.chartData.map(obj => { return obj.count })

      const len = this.chartData.length
      title.length = len
      data.length = len

      const tmp = this.chartData.map(obj => { return obj.succScale })
      const rateData = this.chartData.map(obj => { return obj.succScale })

      rateData.pop()

      const rate = rateData.map((v, i) => {
        const item = {
          value: 0,
          label: {
            show: true,
            formatter: '{a|' + v + '%}'
          }
        }
        return item
      })
      option = {
        dataZoom: [
          {
            show: true,
            realtime: true,
            start: 0,
            bottom: '-8',
            left: 'center',
            end: 50,
            textStyle: false
          },
          {
            textStyle: false,
            type: 'inside',
            realtime: true,
            bottom: '-10',
            left: 'center',
            start: 0,
            end: 50
          }
          /* {
              show: true,
              realtime: true,
              start: 0,
              end: 50
            },
            {
              type: 'inside',
              realtime: true,
              start: 0,
              end: 50
            }*/
        ],
        tooltip: {
          trigger: 'axis',
          axisPointer: { // 坐标轴指示器，坐标轴触发有效
            type: 'shadow' // 默认为直线，可选为：'line' | 'shadow'
          }
        },
        toolbox: {
          feature: {
            saveAsImage: {}
          }
        },
        legend: {
          data: ['漏斗步骤']
        },
        grid: {
          left: '3%',
          right: '4%',
          bottom: '50px',
          containLabel: true
        },
        yAxis: [{
          type: 'value'
        }],
        xAxis: [{
          type: 'category',
          axisTick: {
            show: false
          },
          data: title
        }],
        series: [{
          name: '漏斗步骤',
          type: 'bar',
          showBackground: true,

          backgroundStyle: {
            color: 'rgba(110, 193, 244, 0.2)'
          },
          barCategoryGap: 180,
          label: {
            show: true,
            position: 'inside'
          },
          itemStyle: {
            color: '#6b96f3'
          },
          data: data
        },
        {
          name: '',
          type: 'bar',
          barGap: '-100%',
          label: {
            position: 'right',
            offset: [2, -50],
            formatter: '{a| {c}%}',
            rich: {
              a: {
                align: 'center',
                color: '#000',
                backgroundColor: {
                  image: toPng
                },
                height: 40,
                width: 65,
                fontSize: 12,
                padding: [0, 0]
              }
            }
          },
          data: rate,
          z: 2
        }
        ]
      }
      this.chart.setOption(
        option
      )
    }
  }
}
</script>
