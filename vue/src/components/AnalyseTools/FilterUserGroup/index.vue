<template>
  <div>
    <div
      style="line-height: 18px;font-weight: 500; font-size: 13px; padding: 10px 16px 10px;font-weight: bolder"
    >
      全局筛选用户分群
    </div>
    <div style="padding: 10px">
      <el-select v-model="selectVal" multiple reserve-keyword collapse-tags placeholder="筛选用户分群" clearable filterable size="mini" @change="onchange">
        <el-option
          v-for="(v,k,index) in opt"
          :key="index"
          :label="v.group_name"
          :value="v.id"
        />
      </el-select>
    </div>
  </div>
</template>

<script>

import { UserGroupSelect } from '@/api/user-group'
export default {
  name: 'Index',
  props: {
    value: {
      type: Array,
      default: []
    }
  },
  data() {
    return {
      selectVal: this.value,
      opt: []
    }
  },
  mounted() {
    this.init()
  },
  methods: {
    onchange() {
      this.$emit('input', this.selectVal)
    },
    async init() {
      this.opt = []
      const form = { 'appid': this.$store.state.baseData.EsConnectID }
      const res = await UserGroupSelect(form)
      if (res.code != 0) {
        this.$message({
          offset: 60,
          type: 'error',
          message: res.msg
        })
        return
      }

      if (res.data == null) {
        res.data = []
      }
      this.opt = res.data
    }
  }

}
</script>

<style scoped>

</style>
