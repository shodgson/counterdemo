<script lang="ts">
import { defineComponent } from 'vue'
import NavBar from "@/components/NavBar.vue"
import MessageBar from "@/components/MessageBar.vue"
export default defineComponent({
  name: "Home",
  components: {
    NavBar,
  },
});
</script>

<template>
  <div>
    <message-bar></message-bar>
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
