<template>
  <div :id="className" :class="className" :style="{height:height,width:width}" />
</template>

<script>
// import echarts from 'echarts'
// require('echarts/theme/macarons') // echarts theme
import resize from './mixins/resize'
import { setNotopt } from './mixins/utils'
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

      const eventArr = []
      const targetArr = []
      const eventSet = new Map()
      for (const k in this.chartData) {
        const traceCharts = this.chartData[k]

        eventSet.set(traceCharts['event'][0], 1)
        eventSet.set(traceCharts['event'][1], 1)

        targetArr.push({
          source: traceCharts['event'][0],
          target: traceCharts['event'][1],
          value: traceCharts['sum_user_count']
        })
      }
      eventSet.forEach((v, k, tmp) => {
        const obj = {
          name: k
        }
        eventArr.push(obj)
      })

      const option = {
        toolbox: {
          feature: {
            saveAsImage: {}
          }
        },
        tooltip: {
          trigger: 'item',
          triggerOn: 'mousemove'
        },
        series: [
          {
            type: 'sankey',
            emphasis: {
              focus: 'adjacency'
            },
            nodeAlign: 'left',
            data: eventArr,
            links: targetArr,
            lineStyle: {
              color: 'source',
              curveness: 0.5
            }
          }
        ]
      }
      try {
        this.chart.setOption(
          option
        )
      } catch (e) {
        this.setOptions2()
      }
    },
    setOptions2() {
      const eventArr = []
      const targetArr = []
      const eventSet = new Map()
      for (const k in this.chartData) {
        const traceCharts = this.chartData[k]

        eventSet.set(traceCharts['event'][0] + ' ', 1)
        eventSet.set(traceCharts['event'][1] + '  ', 1)

        targetArr.push({
          source: traceCharts['event'][0] + ' ',
          target: traceCharts['event'][1] + '  ',
          value: traceCharts['sum_user_count']
        })
      }
      eventSet.forEach((v, k, tmp) => {
        const obj = {
          name: k
        }
        eventArr.push(obj)
      })

      const option = {
        toolbox: {
          feature: {
            saveAsImage: {}
          }
        },
        tooltip: {
          trigger: 'item',
          triggerOn: 'mousemove'
        },
        series: [
          {
            type: 'sankey',
            emphasis: {
              focus: 'adjacency'
            },
            nodeAlign: 'left',
            data: eventArr,
            links: targetArr,
            lineStyle: {
              color: 'source',
              curveness: 0.5
            }
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
