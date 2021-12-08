<script lang="ts">
import { defineComponent } from 'vue'
import NavBar from "@/components/NavBar.vue"
export default defineComponent({
  name: "Home",
  components: {
    NavBar,
  },
});
</script>

<template>
  <div>
    <nav-bar></nav-bar>
    <router-view></router-view>
  </div>
</template>
