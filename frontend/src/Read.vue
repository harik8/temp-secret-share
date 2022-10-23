<style>
    .header {
      padding: 60px;
      text-align: center;
    }
</style>
    
<template>
  <div class="container">
    <div class="header">
      <h2>Secret Share</h2>
    </div>
    <form class="form-horizontal" @submit.prevent="getSecret">
      <div class="form-group">
        <label class="control-label col-sm-2" for="secretkey">Secret Key</label>
        <div class="col-sm-10">          
          <input type="secretkey" class="form-control" id="secretkey" placeholder="Enter your Secret Key here" name="secretkey" v-model="secretkey" required>
        </div>
      </div>
      <br>
      <div class="form-group">        
        <div class="col-sm-offset-2 col-sm-10">
          <button type="submit" class="btn btn-primary">Submit</button>
        </div>
      </div>
      <br>
      <div v-if="toggle == true" class="alert alert-primary" style="width:1080px" role="alert">
        <b>{{response}}</b>
      </div>
    </form>
  </div>
</template>

<script>
    import axios from "axios";
    let url = 'https://mmvoql3tvybb2tpv7y7zunxva40rgqxi.lambda-url.us-east-1.on.aws/secret'
    let headers = {
      headers: {
        'Content-Type': 'application/json',
      } 
    }
    export default {
      name: "App",
      data() {
        return {
          response: "",
          toggle: false
        };
    },
    methods: {
        async getSecret() {
          console.log("SID", this.$route.params['secretid'])
          const res = await axios.get(url, { params: 
          {
            SecretID: this.$route.params['secretid'],
            SecretKey: this.secretkey
          }}, headers);
          this.toggle = true
          this.response = res.data;
          this.secretkey = "";
        },
      },
    };
  </script>