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
    <form class="form-horizontal" @submit.prevent="addSecret">
      <div class="form-group">
        <label class="control-label col-sm-20">Message</label>
        <div class="col-sm-10">
          <textarea class="form-control" rows="5" id="message" placeholder="Enter message"  v-model="message" required></textarea>
        </div>
      </div>
      <br>
      <div class="form-group">
        <label class="control-label col-sm-2" for="secretkey">Secret Key</label>
        <div class="col-sm-10">          
          <input type="secretkey" class="form-control" id="secretkey" placeholder="Enter Secret Key" name="secretkey" v-model="secretkey" required>
        </div>
      </div>
      <br>
      <div class="form-group" style="width:1080px; height:100px">
        <label for="expiration">Expire After</label>
          <select class="form-control" id="expiration" placeholder="Select" v-model="activeduration">
            <option value="1h">1 Hour</option>
            <option value="2h">2 Hours</option>
            <option value="12h">12 Hours</option>
            <option value="24h">1 Day</option>
          </select>
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
  export default {
    name: "App",
    data() {
      return {
        response: "",
        toggle: false
      };
    },
    methods: {
      async addSecret() {
        let url = import.meta.env.VITE_ADD_LAMBDA_END_POINT
        let headers = {
          headers: {
            'Content-Type': 'application/json',
          } 
        }

        const res = await axios.post(url, {
          Message: this.message,
          SecretKey: this.secretkey,
          ActiveDuration: this.activeduration
        }, headers);

        this.toggle = true
        // this.response = "Please Share the, Secret Link - [https://main.d1t0nhgavad3l.amplifyapp.com/secret/" + res.data + "] & Secret Key - [" + this.secretkey + "]";
        this.response = "Please Share the, Secret Link - [" + import.meta.env.VITE_WEB_APP_URL + "/secret/" + res.data + "] & Secret Key - [" + this.secretkey + "]";
        this.message = "";
        this.secretkey = "";
        this.activeduration = "";
      },
    },
  };
</script>