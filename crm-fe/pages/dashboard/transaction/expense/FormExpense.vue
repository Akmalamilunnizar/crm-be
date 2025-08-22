<script setup lang="ts">
import { object, string, type InferType } from "yup";
import type { FormSubmitEvent } from "#ui/types";
import { transactionAdminApi } from "@/api/admin/transaction";
import { accountAdminApi } from "@/api/admin/account";

const props = defineProps<{
  id?: string;
  type?: {
    type?: string;
    type_cash?: string;
  };
}>();
const state = reactive({
  type_cash: "",
  type_in_out: "",
  account_id: "",
  date: "",
  description: "",
  file: "",
  amount: "",
  category: "",
  category_custom: "",
  tags: "",
  payer_id: "",
  method: "",
  ref: "",
});
if (props.id) {
  await transactionAdminApi()
    .getTransaction(props.id)
    .then((response) => {
      Object.assign(state, response.data);
    })
    .catch((error) => {
      console.error("Error fetching companies:", error);
    });
}
const schema = object({
  account_id: string()
    .min(3, "Must be at least 3 characters")
    .required("Required"),
  date: string().required("Date is required"),
  amount: string().required("Amount is required"),
});

state.type_in_out = props.type?.type ?? "";
state.type_cash = props.type?.type_cash ?? "";
const isAdvanced = ref(false);

// await companyAdminApi().getAllCompanies().then((response) => {
//   companies.value = response.data.map((value: any, index: number) => ({
//     label: value.name,
//     value: value.id
//   }))
// })

type Schema = InferType<typeof schema>;
async function onSubmit(event: FormSubmitEvent<Schema>) {
  // Do something with event.data
  console.log(event.data);
  if (state.category_custom != "") {
    state.category = state.category_custom
  }
  await transactionAdminApi()
    .createTransactions(state)
    .then((response) => {
      console.log("Success creating / editing company", response);
      onSuccess();
    })
    .catch((error) => {
      console.error("Error creating company:", error);
    });
  // Do something with event.data
  console.log(event.data);
}

const emit = defineEmits(["success"]);

function onSuccess() {
  emit("success");
}
state.date =
  new Date().getFullYear() +
  "-" +
  (new Date().getMonth() + 1).toString().padStart(2, "0") +
  "-" +
  new Date().getDate().toString().padStart(2, "0");
const accounts = ref<any[]>([]);
async function getAccountData() {
  accountAdminApi()
    .getAllAccount()
    .then((response) => {
      accounts.value = response.data.map((value: any, index: number) => ({
        label: value.name,
        value: value.id,
      }));
    });
}
await getAccountData();

const categoryOptions = [
  {
    label: "pemberian pinjaman kasbon",
    value: "pemberian pinjaman kasbon",
  },
  {
    label: "cost of goods sold",
    value: "cost of goods sold",
  },
  {
    label: "salary",
    value: "salary",
  },
  {
    label: "payrol taxes",
    value: "payrol taxes",
  },
  {
    label: "electricity",
    value: "electricity",
  },
  {
    label: "other",
    value: "other",
  },
]

const methodOptions = [
  {
    label: "transfer",
    value: "transfer",
  },
  {
    label: "cash",
    value: "cash",
  },
]
</script>

<template>
  <UModal>
    <div class="p-2 mb-4 text-2xl font-bold text-center">
      <h1>Add Expense</h1>
    </div>
    <div class="p-4">
      <UForm :schema="schema" :state="state" class="space-y-4" @submit="onSubmit">
        <div class="flex gap-4 flex-row-2">
          <div class="w-full">
            <UFormGroup label="Account" name="account">
              <USelectMenu v-model="state.account_id" :options="accounts" value-attribute="value"
                option-attribute="label" />
            </UFormGroup>
            <UFormGroup label="Date" name="date">
              <UInput v-model="state.date" type="date" />
            </UFormGroup>
            <UFormGroup label="Description" name="description">
              <UInput v-model="state.description" />
            </UFormGroup>
            <UFormGroup label="Amount" name="amount">
              <UInput v-model="state.amount" type="number" />
            </UFormGroup>
            <UFormGroup label="Category" name="category">
              <USelectMenu v-model="state.category" :options="categoryOptions" value-attribute="value"
                option-attribute="label" />
            </UFormGroup>
            <UFormGroup v-if="state.category == 'other'" label="Category" name="category">
              <UInput v-model="state.category_custom" />
            </UFormGroup>
            <UFormGroup label="Method" name="method">
              <USelectMenu v-model="state.method" :options="methodOptions" value-attribute="value"
                option-attribute="label" />
            </UFormGroup>
            <!-- <div class="flex justify-end">
              <UButton
                variant="link"
                color="cyan"
                @click="isAdvanced = !isAdvanced"
              >
                Advanced
              </UButton>
            </div> -->
          </div>
          <div class="w-full" v-if="isAdvanced">
            <UFormGroup label="Category" name="category">
              <USelectMenu v-model="state.category" :options="accounts" value-attribute="value"
                option-attribute="label" />
            </UFormGroup>
            <UFormGroup label="Tags" name="tags">
              <UInput v-model="state.tags" />
            </UFormGroup>
            <UFormGroup label="Payer" name="payer">
              <USelectMenu v-model="state.payer_id" :options="accounts" value-attribute="value"
                option-attribute="label" />
            </UFormGroup>
            <UFormGroup label="Method" name="method">
              <USelectMenu v-model="state.method" :options="accounts" value-attribute="value"
                option-attribute="label" />
            </UFormGroup>
            <UFormGroup label="Ref#" name="ref">
              <UInput v-model="state.ref" />
            </UFormGroup>
            <p class="text-xs text-gray-500">e.g. Transaction ID, Check No.</p>
          </div>
        </div>

        <UButton type="submit">Submit</UButton>
      </UForm>
    </div>
  </UModal>
</template>
