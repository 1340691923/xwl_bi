<template>
  <div class="navbar">
    <div class="hamburger-container" style="margin-top: 5px;color: white;margin-left: 10px">

      <router-link style="color: white" to="/">
        {{ title }}
      </router-link>
    </div>
    <div class="hamburger-container" style="margin-top: 5px;color: white;margin-left: 10px">
      <a-divider type="vertical" />
    </div>
    <select-link class="hamburger-container" style="margin-top: 5px" />
    <head-navbar />
    <div class="right-menu">
      <template v-if="device!=='mobile'">
        <el-tooltip content="刷新页面" effect="dark" placement="bottom">
          <reload id="reload" style="font-size:25px" class="right-menu-item hover-effect top_module" />
        </el-tooltip>
        <el-tooltip content="快捷搜索菜单路由" effect="dark" placement="bottom">
          <search id="header-search" class="right-menu-item hover-effect top_module" />
        </el-tooltip>
        <el-tooltip content="是否全屏" effect="dark" placement="bottom">
          <screenfull id="screenfull" class="right-menu-item hover-effect top_module" />
        </el-tooltip>

        <!--<el-tooltip content="全局字体大小" effect="dark" placement="bottom">
        <size-select id="size-select" class="right-menu-item hover-effect top_module"/>
      </el-tooltip>-->
      </template>

      <el-dropdown class="avatar-container right-menu-item hover-effect" trigger="click">
        <!--<img :src="logo" class="user-avatar">-->
        <div class="avatar-wrapper" style="background: white;height: 40px;line-height: 40px;width: 40px;text-align: center;border-radius: 50%;cursor: pointer;">
          {{ name.substr(0,1) }}
        </div>
        <el-dropdown-menu slot="dropdown">
          <router-link to="/">
            <el-dropdown-item>首页</el-dropdown-item>
          </router-link>
          <el-dropdown-item divided>
            <span style="display: block;" @click="dialogVisible = true">修改密码</span>
          </el-dropdown-item>
          <el-dropdown-item divided>
            <span style="display: block;" @click="logout">注销</span>
          </el-dropdown-item>
        </el-dropdown-menu>
      </el-dropdown>
      <el-dialog
        :close-on-click-modal="false"
        :visible.sync="dialogVisible"
        title="修改密码"
      >
        <el-form label-width="100px" label-position="left">
          <el-form-item label="密码">
            <el-input v-model="password" style="width: 300px" placeholder="密码" />
          </el-form-item>
          <el-form-item label="再次输入密码">
            <el-input v-model="password2" style="width: 300px" placeholder="再次输入密码" />
          </el-form-item>
        </el-form>
        <div style="text-align:right;">
          <el-button type="danger" icon="el-icon-close" @click="dialogVisible=false">取消</el-button>
          <el-button type="primary" icon="el-icon-check" @click="modifyPass">确认</el-button>
        </div>
      </el-dialog>
    </div>
  </div>
</template>

<script>
import logo from '@/assets/index.ico'
import Breadcrumb from '@/components/Breadcrumb'
import Hamburger from '@/components/Hamburger'
import Screenfull from '@/components/Screenfull'
// import SizeSelect from '@/components/SizeSelect'
import Search from '@/components/HeaderSearch'
import SelectLink from '@/components/SelectLink'
import HeadNavbar from '../components/HeadNavbar'
import Reload from '@/components/Reload'
import { ModifyPwd } from '@/api/user'
import { mapGetters } from 'vuex'

export default {
  components: {
    Breadcrumb,
    Hamburger,
    Reload,
    Screenfull,
    // SizeSelect,
    Search,
    SelectLink,
    HeadNavbar
  },
  data() {
    return {
      logo: logo,
      password: '',
      password2: '',
      dialogVisible: false
    }
  },
  computed: {
    ...Vuex.mapGetters([
      'sidebar',
      'avatar',
      'device',
      'name'
    ]),
    title() {
      return process.env.VUE_APP_BASE_TITLE
    }
  },
  methods: {
    async modifyPass() {
      this.password2 = this.password2.trim()
      this.password = this.password.trim()
      if (this.password == '' || this.password2 == '') {
        this.$message({
          offset: 60,
          type: 'error',
          message: '密码不能为空'
        })
        return
      }
      if (this.password2 != this.password) {
        this.$message({
          offset: 60,
          type: 'error',
          message: '先后两次密码不一致'
        })
        return
      }
      this.$confirm('确定修改密码吗?', '警告', {
        confirmButtonText: '确认',
        cancelButtonText: '取消',
        type: 'warning'
      })
        .then(async() => {
          const res = await ModifyPwd({ password: this.password })
          this.password2 = ''
          this.password = ''
          if (res.code != 0) {
            this.$message({
              offset: 60,
              type: 'error',
              message: res.msg
            })
            return
          }
          this.$message({
            offset: 60,
            type: 'success',
            message: res.msg
          })
          this.dialogVisible = false
        })
        .catch(err => {
          console.error(err)
        })
    },
    toggleSideBar() {
      this.$store.dispatch('app/toggleSideBar')
    },
    async logout() {
      await this.$store.dispatch('user/logout')
      this.$router.push(`/login?redirect=${this.$route.fullPath}`)
    }
  }
}
</script>

<style lang="scss" scoped>
  .navbar {
    height: 55px;
    overflow: hidden;
    position: relative;
    background: #475285;
    box-shadow: 0 1px 4px rgba(0, 21, 41, .08);

  .hamburger-container {

    line-height: 46px;
    height: 100%;
    float: left;
    cursor: pointer;
    transition: background .3s;
    -webkit-tap-highlight-color: transparent;

  &
  :hover {
    background: rgba(0, 0, 0, .025)
  }

  }

  .breadcrumb-container {
    float: left;
  }

  .errLog-container {
    display: inline-block;
    vertical-align: top;
  }

  .right-menu {
    float: right;
    height: 100%;
    line-height: 50px;

  &
  :focus {
    outline: none;
  }

  .right-menu-item {
    display: inline-block;
    padding: 0 8px;
    height: 100%;
    font-size: 18px;
    color: #5a5e66;
    vertical-align: text-bottom;

  &
  .hover-effect {
    cursor: pointer;
    transition: background .3s;

  &
  :hover {
    background: rgba(0, 0, 0, .025)
  }

  }
  }

  .avatar-container {
    margin-right: 30px;

  .avatar-wrapper {
    margin-top: 5px;
    position: relative;

  .user-avatar {
    cursor: pointer;
    width: 40px;
    height: 40px;
    border-radius: 10px;
  }

  .el-icon-caret-bottom {
    cursor: pointer;
    position: absolute;
    right: -20px;
    top: 25px;
    font-size: 12px;
  }

  }
  }
  }
  }
  .all-fontsize {
    font-size: 15px !important;
  }

  .ta-avatar-icon-name-icon {
    display: inline-block;
    align-items: center;
    justify-content: center;
    width: 24px;
    height: 24px;
    font-size: 12px;
    line-height: 22px;
    text-align: center;
    border-radius: 20px;
  }

  .avatar___hTjhs {
    height: 32px;
    margin: 0 12px;
    padding: 0 12px;
    line-height: 32px;
    cursor: pointer;
  }
  .top_module:hover{
    cursor: pointer;
  }

</style>
