<template>
  <div id="app" class="container">
    <div class="row">
      <div class="col">
        <ping remote="https://htping-a.vis.ethz.ch" name="old"></ping>
      </div>
      <div class="col">
        <ping remote="https://htping-b.vis.ethz.ch" name="new"></ping>
      </div>
    </div>
  </div>
</template>

<script>
import Ping from './components/Ping.vue'

export default {
  name: 'app',
  components: {
    Ping
  },
}
</script>

<style>
#app {
  font-family: 'Avenir', Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  text-align: center;
  color: #2c3e50;
  margin-top: 60px;
}
</style>
