new Vue({
  el: '#app',
  data: {
    pageNo: 1,
    name_part: "",
    orders: [],
    order_date_from: null,
    order_date_to: null,
    API_URL: config.API_URL,
    TZ_STR: "Australia/Melbourne|AEST AEDT|-a0 -b0|0101010101010101010101010101010101010101010101010101010101010101010101010101010101010101010101010101010101010101010101010101010101010101010101|-293lX xcX 10jd0 yL0 1cN0 1cL0 1fB0 19X0 17c10 LA0 1C00 Oo0 1zc0 Oo0 1zc0 Oo0 1zc0 Rc0 1zc0 Oo0 1zc0 Oo0 1zc0 Oo0 1zc0 Oo0 1zc0 Oo0 1zc0 Rc0 1zc0 Oo0 1zc0 Oo0 1zc0 Oo0 1zc0 U00 1qM0 WM0 1qM0 11A0 1tA0 U00 1tA0 U00 1tA0 Oo0 1zc0 Oo0 1zc0 Rc0 1zc0 Oo0 1zc0 WM0 1qM0 11A0 1o00 11A0 1o00 11A0 1o00 11A0 1qM0 WM0 11A0 1o00 1qM0 11A0 1o00 11A0 1o00 11A0 1qM0 WM0 1qM0 11A0 1o00 WM0 1qM0 14o0 1cM0 1cM0 1cM0 1cM0 1cM0 1cM0 1cM0 1cM0 1fA0 1cM0 1cM0 1cM0 1cM0 1cM0 1cM0 1cM0 1cM0 1cM0 1cM0 1cM0 1fA0 1cM0 1cM0 1cM0 1cM0 1cM0 1cM0 1cM0 1cM0 1cM0 1cM0 1fA0 1cM0 1cM0 1cM0 1cM0 1cM0 1cM0 1cM0 1cM0 1cM0 1cM0 1fA0 1cM0 1cM0 1cM0 1cM0 1cM0 1cM0 1cM0 1cM0 1cM0 1cM0 1cM0 1fA0 1cM0 1cM0 1cM0 1cM0|39e5",
    PAGE_SIZE: 5,
  },
  computed: {
    query: function() {
      var date_from = this.order_date_from;
      var date_to = this.order_date_to;
      return {
        DateFrom: date_from ? moment(date_from).utc().format() : "",
        DateTill: date_to ? moment(date_to).add(1, "days").utc().format() : "",
        PartOfName: this.name_part,
        PageNo: this.pageNo,
        PageSize: this.PAGE_SIZE,
      }      
    },
  },
  watch: {
    order_date_from: function(value) {
      this.search_orders(this.query);
    },
    order_date_to: function(value) {
      this.search_orders(this.query);
    },
    pageNo:  function(value, oldValue) {
      if (value < 1) {
        this.pageNo = oldValue;
      }
      this.search_orders(this.query, oldValue);
    },
  },
  methods: {
    clear_query: function() {
      this.name_part = "";
      this.order_date_from = null;
      this.order_date_to = null;
      this.pageNo = 1;
    },
    transform_tz: function(date_str) {
      return moment(date_str).tz("Australia/Melbourne").format("MMMM Do YYYY, h:mm a");
    },
    search_orders: function(query, oldPage) {
      var vm = this;
      $.ajax({
        url: vm.API_URL + "/api/orders/search",
        data: query,
        success: function(data, textStatus, xhr) {
          vm.orders = data;
          if (oldPage && vm.orders.length === 0 && vm.pageNo > 1) {
            vm.pageNo = oldPage;
          }
        },
        complete: function(xhr, textStatus) {
          if (xhr.status !== 200) {
            console.log("status code: " + xhr.status);
            console.log("text status: " + textStatus);
          }
        }
      });
    },
  },
  mounted: function() {
    this.search_orders({});
    moment.tz.add(this.TZ_STR);
  },
});
