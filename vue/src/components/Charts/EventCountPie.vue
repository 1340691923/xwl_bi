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
      const option = {
        legend: {
          top: 'bottom'
        },
        grid: {
          left: '3%',
          right: '50%',
          bottom: '50px',
          containLabel: true
        },
        toolbox: {
          show: true,
          feature: {
            saveAsImage: { show: true }
          },
          trigger: 'item',
          confine: true
        },
        series: [
          {
            label: {
              position: 'outer',
              alignTo: 'edge',
              margin: 10
              /* normal: {
                  position: "inside"//此处将展示的文字在内部展示
                }*/
            },
            name: '事件所占比例图',
            type: 'pie',
            radius: ['10%', '30%'],
            center: ['50%', '50%'],
            roseType: 'area',
            itemStyle: {
              normal: {
                label: {
                  show: true,
                  formatter: '{b} : {c} ({d}%)'
                },
                labelLine: { show: true }
              },
              borderRadius: 8
            },
            data: this.chartData
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
