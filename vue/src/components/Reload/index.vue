<template>
  <div>
    <i class="el-icon-refresh" style="color:#ffd04b" @click="click" />
  </div>
</template>

<script>

export default {
  name: 'Reload',
  inject: ['reload'],
  data() {
    return {

    }
  },
  mounted() {

  },
  beforeDestroy() {

  },
  methods: {
    click() {
      this.reload()
    }
  }
}
</script>

<style scoped>

</style>
