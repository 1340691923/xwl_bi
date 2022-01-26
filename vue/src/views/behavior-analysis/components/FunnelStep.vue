<template>
  <a-cascader
    v-model="step"
    style="width: 160px"
    :options="stepOptions"
    placeholder="请选择步骤"
    @change="changeStep"
  />
</template>

<script>
export default {
  name: 'FunnelStep',
  props: {
    value: {
      type: Number,
      default: 0
    },
    funnelRes: {
      type: Array,
      default: []
    },
    groupData: {
      type: Array,
      default: []
    },
    tableHeader: {
      type: Array,
      default: []
    }
  },
  data() {
    return {
      step: this.value,
      stepOptions: [],
      groupDataShow: [],
      tableHeaderShow: [],
      showList: [],
      tableData: []
    }
  },
  watch: {
    step(newV, oldV) {
      this.$emit('input', this.step)
    }
  },
  methods: {
    changeStep() {
      switch (this.step.length) {
        case 0:
          this.step = [-1]
        case 1:
          this.showList = JSON.parse(JSON.stringify(this.funnelRes))
          this.groupDataShow = JSON.parse(JSON.stringify(this.groupData))
          this.tableHeaderShow = JSON.parse(JSON.stringify(this.tableHeader))

          break
        case 2:
          const stepStart = this.step[0]
          const stepOver = this.step[1]

          const funnelRes = []

          for (const v of this.funnelRes) {
            if (v.level_index >= stepStart && v.level_index <= stepOver) {
              funnelRes.push(v)
            }
          }
          const groupDataShow = {}
          for (const k in this.groupData) {
            const tmp = []

            for (const v of this.groupData[k]) {
              if (v.level_index >= stepStart && v.level_index <= stepOver) {
                tmp.push(v)
              }
            }
            groupDataShow[k] = tmp
          }
          this.groupDataShow = groupDataShow
          this.showList = funnelRes
          const groupTitle = this.tableHeader[0]
          const tableHeader = [groupTitle]
          for (const k in this.tableHeader) {
            if (k >= stepStart && k <= stepOver) {
              tableHeader.push(this.tableHeader[k])
            }
          }

          this.tableHeaderShow = tableHeader
      }
      this.changeTable()
      this.refreshData()
    },

    changeTable() {
      const tableData = []

      for (const k in this.groupDataShow) {
        const countArr = []
        for (const v of this.groupDataShow[k]) {
          countArr.push({
            count: v.count,
            conversionScale: this.NaN2Zero(v.conversionScale),
            washScale: this.NaN2Zero(v.washScale)
          })
        }
        tableData.push({ 'groupKey': k, countArr: countArr })
      }

      this.tableData = tableData
    },
    init() {
      const stepOpt = [
        {
          value: -1,
          label: '全步骤'
        }
      ]
      for (const i in this.funnelRes) {
        if (i == this.funnelRes.length - 1) {
          break
        }

        const stepNum = Number(i) + 1
        var obj = {
          value: stepNum,
          label: '步骤' + stepNum,
          children: []
        }
        let childrenStep = Number(stepNum) + 1

        for (const k in [...new Array(this.funnelRes.length - stepNum).keys()]) {
          obj.children.push({
            value: childrenStep,
            label: '步骤' + childrenStep
          })
          childrenStep++
        }
        stepOpt.push(obj)
      }
      this.stepOptions = stepOpt
    }
  }

}
</script>

<style scoped>

</style>
