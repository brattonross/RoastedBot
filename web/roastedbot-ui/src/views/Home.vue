<template>
  <v-container fluid grid-list-md>
    <v-layout row wrap>
      <v-flex
        xs12
        v-for="channel in channels"
        :key="channel.name">
        <v-card>
            <v-card-text>{{ channel.name }}</v-card-text>
        </v-card>
      </v-flex>
    </v-layout>
  </v-container>
</template>

<script>
import axios from 'axios'

export default {
  data: () => ({
    channels: []
  }),
  mounted () {
    axios
      .get('http://localhost:9001/channels')
      .then(resp => {
        this.channels = resp.data
      })
      .catch(console.error)
  }
}
</script>
