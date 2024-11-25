<template>
  <el-form :model="settingRef">
    <el-form-item label="package前缀">
      <el-input v-model="settingRef.package_prefix"/>
    </el-form-item>
    <el-form-item label="路径前缀">
      <el-input v-model="settingRef.target"/>
    </el-form-item>
        <el-form-item label="前端路径前缀">
          <el-input v-model="settingRef.web_target"/>
        </el-form-item>

  </el-form>
  <el-button @click="load">加载</el-button>
  <el-button @click="save">保存</el-button>
</template>


<script setup lang="ts">
import {reactive} from "vue";
import {loadGoInfo} from "@/api/module/setting";

const settingRef = reactive({
  target: '',
  web_target:'',
  package_prefix: ''
})

init()

async function load() {
  settingRef.package_prefix = await loadGoInfo(settingRef.target)
}

function init() {
  settingRef.target = localStorage.getItem("go_target")!!
  settingRef.web_target = localStorage.getItem("web_target")!!
  settingRef.package_prefix = localStorage.getItem("package_prefix")!!
}

function save() {
  localStorage.setItem("package_prefix", settingRef.package_prefix)
  localStorage.setItem("go_target", settingRef.target)
  localStorage.setItem("web_target", settingRef.web_target)
}

</script>


<style scoped>

</style>
