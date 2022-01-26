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
      const title = this.chartData.map(obj => { return obj.date_group })
      const data = this.chartData.map(obj => { return obj.count })

      const option = {
        tooltip: {
          trigger: 'axis',
          axisPointer: {
            type: 'shadow'
          }
        },
        grid: {
          left: '3%',
          right: '4%',
          bottom: '3%',
          containLabel: true
        },
        xAxis: [
          {
            type: 'category',
            data: title,
            axisTick: {
              alignWithLabel: true
            }
          }
        ],
        yAxis: [
          {
            type: 'value'
          }
        ],
        series: [
          {
            name: '事件触发数',
            type: 'bar',
            barWidth: '60%',
            data: data
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
