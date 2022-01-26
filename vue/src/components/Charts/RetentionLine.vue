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
    xData: {
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

      const legendData = []
      const xAxisData = this.xData
      const seriesData = []
      for (const v of this.chartData) {
        legendData.push(v['dates'])

        const obj = {
          name: v['dates'],
          type: 'line',
          smooth: true,
          animationDuration: 2800,
          animationEasing: 'quadraticOut'
        }

        if (this.isScale) {
          obj['data'] = v['conversionScaleArr'].slice(1, v['conversionScaleArr'].length)
        } else {
          obj['data'] = v['value'].slice(1, v['value'].length)
        }
        seriesData.push(obj)
      }

      option = {
        tooltip: {
          trigger: 'axis',
          confine: true
        },
        legend: {
          top: '93%',
          data: legendData,
          type: 'scroll',
          selector: ['all', 'inverse']
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
