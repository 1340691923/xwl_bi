<template>

  <a-tooltip placement="top" style="cursor: pointer">
    <template slot="title">
      <span>创建用户分群</span>
    </template>
    <a-button
      type="link"
      icon="usergroup-add"
      style="margin-left: 4px;"
      class="right_icon"
      @click="add"
    />
  </a-tooltip>
</template>

<script>
import { AddUserGroup } from '@/api/user-group'

export default {
  name: 'AddUserGroup',
  props: {
    uid: {
      type: Array,
      default: []
    }
  },
  methods: {
    add() {
      const uid = this.uid
      this.$prompt('请输入用户分群名', '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消'
      }).then(async({ value }) => {
        const form = { 'appid': this.$store.state.baseData.EsConnectID, uids: uid, name: value }

        const res = await AddUserGroup(form)

        if (res.code != 0) {
          this.$message({
            type: 'error',
            message: res.msg
          })
          return
        }
        this.$message({
          type: 'success',
          message: res.msg
        })
      }).catch(err => {
        console.log(err)
      })
    }
  }

}
</script>

<style scoped>

</style>
