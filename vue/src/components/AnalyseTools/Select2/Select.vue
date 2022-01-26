<template>
  <el-select
    v-model="selectedArray"
    size="mini"
    multiple
    collapse-tags
    :placeholder="placeholder"
    @change="changeSelect"
  >
    <el-checkbox v-model="checked" @change="selectAll">全选</el-checkbox>
    <el-option v-for="(item, index) in options" :key="index" :label="item.label" :value="item.value">
      <span style="float: left">{{ item.label }}</span>
      <span style="float: right; color: #8492a6; font-size: 13px">{{ item.value }}</span>
    </el-option>
  </el-select>
</template>

<script>
export default {
  name: 'Select',
  props: {
    value: {
      type: Array,
      default: []
    },
    options: {
      type: Array,
      default: []
    },
    placeholder: {
      type: String,
      default: ''
    },
    checkeds: {
      type: Boolean,
      default: false
    }
  },
  data() {
    return {
      checked: this.checkeds,
      selectedArray: this.value
    }
  },
  methods: {
    selectAll() {
      this.selectedArray = []
      if (this.checked) {
        this.options.map((item) => {
          this.selectedArray.push(item.value)
        })
      } else {
        this.selectedArray = []
      }
      this.$emit('input', this.selectedArray)
    },
    changeSelect(val) {
      if (val.length === this.options.length) {
        this.checked = true
      } else {
        this.checked = false
      }
      this.$emit('input', val)
    }
  }

}
</script>

<style scoped>
.el-checkbox {
  text-align: right;
  width: 100%;
  padding-right: 10px;
}
</style>
