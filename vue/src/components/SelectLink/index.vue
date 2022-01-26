<template>
  <div class="header-search ">

    <el-select
      v-model="linkID"
      class="header-search-select"
      filterable
      default-first-option
      placeholder="请选择您的应用"
      @change="change"
    >
      <el-option :value="Number(0)" label="请选择您的应用" />
      <el-option v-for="item in opt" :key="item.id" :value="Number(item.id)" :label="item.app_name" />
    </el-select>

  </div>
</template>

<script>
import { Config } from '@/api/app'

export default {
  inject: ['reload'],
  name: 'SelectLink',
  data() {
    return {
      opt: [],
      linkID: '',
      time: null,
      timeSecend: 60
    }
  },
  computed: {},
  watch: {},
  mounted() {
    const obj = this.$store.state.baseData.EsConnectID
    this.linkID = Number(obj)
    this.getEsOpt()
    // this.startLoop()
  },

  beforeDestroy() {
    // 清除定时器
    clearInterval(this.time)
    this.time = null
  },
  methods: {
    startLoop() {
      this.time = setInterval(() => {
        this.getEsOpt()
      }, this.timeSecend * 1000)
    },
    async getEsOpt() {
      const res = await Config()
      this.opt = res.data.list
    },
    refresh() {
      this.getEsOpt()
      this.$message({
        offset: 60,

        type: 'success',
        message: '刷新应用信息成功'
      })
    },
    change(link) {
      this.$store.dispatch('baseData/SetEsConnect', link)
      this.reload()
    }
  }

}
</script>

<style>
  /deep/ .el-input{
    font-size: 0.32rem !important;
  }
  /deep/ .el-input--suffix .el-input__inner {
    color: rgb(244, 65, 68);
  }
  /deep/ .el-select-dropdown__item.selected {
    color: rgb(244, 65, 68);
  }
</style>

<style lang="scss" scoped>

  .header-search {
    font-size: 0 !important;

    .header-search-select {

      /deep/ .el-input__inner {
        color: white;
        background: #475285;
        border-radius: 0;
        border: 0;
        padding-left: 0;
        padding-right: 0;
        box-shadow: none !important;
        border-bottom: 1px solid #d9d9d9;
        vertical-align: middle;
        margin-left: 10px;
      }
    }

  }

</style>
