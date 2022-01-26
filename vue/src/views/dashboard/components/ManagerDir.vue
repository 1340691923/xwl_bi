<template>
  <div>
    <a-input v-model="panelStr" style="width: 170px" class="search" placeholder="搜索看板" allow-clear @input="change">
      <a-icon slot="prefix" type="search" />
    </a-input>
  </div>
</template>

<script>
export default {
  name: 'ManagerDir',
  props: {
    value: {
      type: String,
      default: ''
    }
  },
  data() {
    return {
      panelStr: this.value
    }
  },
  methods: {
    change() {
      this.$emit('input', this.panelStr)
      this.$emit('change')
    }
  }
}
</script>

<style lang="scss" scoped>
.search {

input {
  border: none !important;
  box-shadow: none !important;
}

}
</style>
