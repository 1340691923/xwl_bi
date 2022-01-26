<template>
  <div :id="className" :class="className" :style="{height:height,width:width}" />
</template>

<script>
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

    isScale: {
      type: Boolean,
      default: false
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

      const xAxisData = []
      const seriesData = []
      const data = []
      for (const v of this.chartData) {
        xAxisData.push(v['dates'])
        if (this.isScale) {
          data.push(v['conversionScaleArr'][v['conversionScaleArr'].length - 1])
        } else {
          data.push(v['value'][v['value'].length - 1])
        }
      }

      const obj = {
        name: '总体',
        type: 'line',
        smooth: true,
        animationDuration: 2800,
        animationEasing: 'quadraticOut'
      }
      obj['data'] = data

      seriesData.push(obj)

      option = {
        tooltip: {
          trigger: 'axis',
          confine: true
        },
        grid: {
          left: '3%',
          right: '4%',
          bottom: '50px',
          containLabel: true
        },
        toolbox: {
          feature: {
            saveAsImage: {}
          }
        },
        xAxis: {
          type: 'category',
          boundaryGap: false,
          data: xAxisData
        },
        yAxis: {
          type: 'value',
          axisLabel: {
            formatter: this.isScale ? '{value}%' : '{value}个用户'
          }
        },
        series: seriesData
      }

      this.chart.setOption(
        option
      )
    }
  }
}
</script>
<style>
  .z_index_first{
    z-index: 999999999;
  }
</style>
