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
    },
    showLabel: {
      type: String,
      default: '数量'
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

        tooltip: {
          trigger: 'item'
        },
        toolbox: {
          feature: {
            saveAsImage: {}
          }
        },
        legend: {
          orient: 'vertical',
          left: 'left'
        },
        series: [
          {
            name: this.showLabel,
            type: 'pie',
            radius: '50%',
            data: this.chartData,
            emphasis: {
              itemStyle: {
                shadowBlur: 10,
                shadowOffsetX: 0,
                shadowColor: 'rgba(0, 0, 0, 0.5)'
              }
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
