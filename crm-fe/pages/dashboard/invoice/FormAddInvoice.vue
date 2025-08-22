<script setup lang="ts">
import { object, string, type InferType } from "yup";
import type { FormSubmitEvent } from "#ui/types";
import { customerAdminApi } from "@/api/admin/customer";
import { invoiceAdminApi } from "@/api/admin/invoice";
import { internetPackageAdminApi } from "@/api/admin/internet-package";

const props = defineProps({
  isEdit: {
    type: Boolean,
    required: false,
  },
  data: {
    type: Object,
    default: () => ({
      customer_id: {
        type: String,
        default: "",
      },
      amount: {
        type: String,
        default: "",
      },
    }),
  },
});
const loadingProduct = ref(false);

const schema = object({
  customer_id: string().required("Customer is required"),
  amount: string().required("Amount is required"),
});

type Schema = InferType<typeof schema>;

const state = reactive({
  customer_id: "",
  amount: 0,
  invoice_items: [
    {
      name: "",
      qty: 0,
      price: 0,
      total: 0,
    },
  ],
});

function addItem() {
  state.invoice_items.push({
    name: "",
    qty: 0,
    price: 0,
    total: 0,
  });
}
function removeItem(index: number) {
  state.invoice_items.splice(index, 1);
}
function updateTotal(index: number) {
  const item = state.invoice_items[index];
  state.invoice_items[index] = {
    ...item,
    total: item.qty * item.price,
  };

  state.amount = state.invoice_items.reduce((acc, item) => acc + item.total, 0);
  console.log(state.invoice_items[index].total);
}

watch(
  () => props.isEdit,
  (newValue) => {
    if (newValue) {
      (state.customer_id = props.data.customer_id),
        (state.amount = props.data.amount);
    }
  },
  { immediate: true }
);

watch(
  () => props.isEdit,
  (newValue) => {
    if (newValue) {
      (state.customer_id = props.data.customer_id),
        (state.amount = props.data.amount);
    }
  },
  { immediate: true }
);
const emit = defineEmits(["success"]);

function onSuccess() {
  emit("success");
}

async function onSubmit(event: FormSubmitEvent<Schema>) {
  // Do something with event.data
  if (props.isEdit) {
    invoiceAdminApi()
      .editInvoice(props.data.id, state)
      .then((response) => {
        useToast().add({ title: response.message });
        onSuccess();
      })
      .catch((error) => { });
  } else {
    invoiceAdminApi()
      .createInvoice(state)
      .then((response) => {
        useToast().add({ title: response.message });
        onSuccess();
      })
      .catch((error) => { });
  }
}

const customer = ref([]);
const searchOptions = ref();
const productOptions = ref<any[]>([])
const productOptionsD = ref<any[]>([])

async function getDataOptions() {
  customerAdminApi()
    .getAllCustomers()
    .then((response) => {
      customer.value = response.data.map((value: any, index: number) => ({
        label: value.email,
        value: value.id,
      }));
    });
}

function search(q: any) {
  const data = productOptions.value.filter((option: any) =>
    option.toLowerCase().includes(q.toLowerCase())
  );
  if (data.length > 0) {
    return data
  }

  return [q]
}

await internetPackageAdminApi()
  .getAllInternetPacket()
  .then((response) => {
    productOptions.value = response.data.map((value: any, index: number) => value.name);
    productOptionsD.value = response.data.map((value: any, index: number) => ({
      id: value.id,
      label: value.name,
      value: value.name,
      price: value.price,
    }));
  })
  .finally(() => {
    loadingProduct.value = false;
  });
await getDataOptions();

function checkProductIsExist(name: string, index: number) {
  const product = productOptionsD.value.find((option: any) => option.label === name);
  if (product) {
    const item = state.invoice_items[index];
    state.invoice_items[index] = {
      ...item,
      price: product.price,
    };
  }
}
</script>

<template>
  <UModal>
    <div class="w-full p-4">
      <div class="p-2 mb-4 text-2xl font-bold text-center">
        <h1>{{ props.isEdit ? "Edit" : "Add New" }} Invoice</h1>
      </div>
      <UForm :schema="schema" :state="state" class="space-y-4" @submit="onSubmit">
        <UFormGroup label="Customer" name="customer_id">
          <USelectMenu v-model="state.customer_id" :options="customer" value-attribute="value"
            option-attribute="label" />
        </UFormGroup>
        <UFormGroup label="Amount" name="amount">
          <UInput v-model="state.amount" type="number" />
        </UFormGroup>
        <div v-for="(item, index) in state.invoice_items" :key="index" class="space-y-4">
          <UFormGroup :label="`Product ${index + 1} Name`" :name="`item-name-${index}`">
            <UInputMenu v-model="item.name" :loading="loadingProduct" by="id" :options="productOptions"
              @change="(name) => checkProductIsExist(name, index)" :search="search" />
          </UFormGroup>

          <div class="flex space-x-4">
            <UFormGroup label="Quantity">
              <UInput v-model.number="item.qty" type="number" @update:modelValue="() => updateTotal(index)" />
            </UFormGroup>
            <UFormGroup label="Price">
              <UInput v-model.number="item.price" type="number" @update:modelValue="() => updateTotal(index)" />
            </UFormGroup>
            <UFormGroup label="Total">
              <UInput v-model.number="item.total" type="number" disabled />
            </UFormGroup>
          </div>

          <UButton color="red" variant="soft" @click="removeItem(index)" v-if="state.invoice_items.length > 1">
            Hapus Item
          </UButton>
        </div>
        <UFormGroup>
          <UButton @click="addItem" variant="outline">Tambah Item</UButton>
        </UFormGroup>
        <UButton type="submit"> Submit </UButton>
      </UForm>
    </div>
  </UModal>
</template>
