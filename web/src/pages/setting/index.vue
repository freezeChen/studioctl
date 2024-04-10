<script setup lang="ts">
import {ref} from "vue";

const settingRef = ref()


</script>

<template>
  <el-form :model="settingRef">
    <el-form-item label="gomod">
      <el-input/>
    </el-form-item>
    <el-form-item label="路径前缀">
      <el-input/>
    </el-form-item>

  </el-form>

</template>

<style scoped>

</style>
