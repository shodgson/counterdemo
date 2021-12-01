<script lang="ts">
import NavBar from "@/components/NavBar.vue"
export default {
  name: "Home",
  components: {
    NavBar,
  },
};
</script>

<template>
  <div>
    <nav-bar></nav-bar>
    <router-view></router-view>
  </div>
</template>

<style>
#app {
  font-family: Avenir, Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  text-align: center;
  color: #2c3e50;
}
</style>
