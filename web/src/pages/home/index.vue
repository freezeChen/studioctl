<template>
  <el-container>
    <el-aside>
      <el-space>
        studio-tool
      </el-space>

      <el-menu
          router
          collapse
          style="width: 100px;height: 100vh"
      >
        <el-menu-item index="gencode">
          页面
        </el-menu-item>
        <el-menu-item index="setting">
          设置
        </el-menu-item>
      </el-menu>
    </el-aside>
    <el-main>
      <router-view/>
    </el-main>
  </el-container>

</template>
