<template>

  <el-cascader
    v-model="selectVal"
    class="topselect w130 ml15"
    size="mini"
    filterable
    :placeholder="placeholder"
    separator="."
    :options="options"
    :props="{ expandTrigger: 'hover' }"
    @change="onChange"
  />

</template>

<script>
export default {
  name: 'CountSelect',
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
    renderFormat({ labels, selectedOptions }) {
      return labels.join('.')
    }
  },
  data() {
    return {
      selectVal: this.value
    }
  },
  watch: {
    value: {
      handler() {
        this.selectVal = this.value
      },
      deep: true
    }
  },
  beforeMount() {
    this.selectVal = this.value
  },
  methods: {
    onChange(a, b, c) {
      this.$emit('input', this.selectVal)
    },
    filter(inputValue, path) {
      return path.some(option => option.label.toLowerCase().indexOf(inputValue.toLowerCase()) > -1)
    }
  }
}
</script>

<style lang="scss"  scoped>
  .topselect {
  ::v-deep {
  .el-input__inner {
    height: 32px;

  }

  .el-input__prefix, .el-input__suffix {
    height: 30px;
  }

  /* 下面设置右侧按钮居中 */
  .el-input__suffix {
    top: 5px;
  }

  .el-input__icon {
    line-height: inherit;
  }

  .el-input__suffix-inner {
    display: inline-block;
  }
  }
  }

</style>
