<template>

  <a-range-picker v-model="filterDate" format="YYYY-MM-DD" :ranges="dataRange" @change="filterDateCall">
    <div style="cursor: pointer">
      <a-button icon="clock-circle">

        <template v-if="moment().diff(filterDate[1], 'day') == 0 && moment().diff(filterDate[0], 'day') == 0">
          今天
        </template>
        <template v-else-if="moment().diff(filterDate[1], 'day') == 0">
          过去{{ moment().add(1, 'days').diff(filterDate[0], 'day') }}天
        </template>
        <template v-else-if="moment().diff(filterDate[1], 'day') == 1 && moment().diff(filterDate[0], 'day') == 1">
          昨天
        </template>
        <template v-else>
          {{ filterDate[0] }}~{{ filterDate[1] }}
        </template>
      </a-button>
    </div>
  </a-range-picker>

</template>

<script>
import moment from 'moment'

const dataRange = {
  '今天': [moment().startOf('day'), moment()],
  '昨天': [
    moment().startOf('day').subtract(1, 'days'),
    moment().startOf('day').subtract(1, 'days')
  ],
  '最近一周': [moment().startOf('day').subtract(1, 'weeks').add(1, 'days'), moment()],
  '最近两周': [moment().startOf('day').subtract(2, 'weeks').add(1, 'days'), moment()],
  '最近1个月': [moment().startOf('day').subtract(1, 'months').add(1, 'days'), moment()],
  '最近3个月': [moment().startOf('day').subtract(3, 'months').add(1, 'days'), moment()],
  '最近半年': [moment().startOf('day').subtract(6, 'months').add(1, 'days'), moment()],
  '最近1年': [moment().startOf('day').subtract(1, 'years').add(1, 'days'), moment()]
}

export default {
  name: 'Date',
  props: {
    value: {
      type: Array,
      default: []
    }
  },
  data() {
    return {
      dataRange,
      filterDate: this.value
    }
  },
  methods: {
    moment,
    filterDateCall(dates, dateStrings) {
      this.filterDate = dateStrings
      this.$emit('changeDate', this.filterDate)
    }
  }
}
</script>

<style scoped>

</style>
