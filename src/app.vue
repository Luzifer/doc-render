<template>
  <div class="container-fluid py-3 h-100">
    <div class="row h-100">
      <div class="col">
        <div class="input-group mb-3">
          <div class="form-floating">
            <select
              id="selectedSet"
              v-model="selectedSet"
              class="form-select"
            >
              <option
                v-for="dt in docTypes"
                :key="dt.name"
                :value="dt.name"
              >
                {{ dt.description }}
              </option>
            </select>
            <label for="selectedSet">Dokumententyp</label>
          </div>
          <button
            class="btn btn-success"
            :disabled="documentLoading || !modelValid"
            @click="render"
          >
            <i :class="`fas fa-fw me-1 ${documentLoading ? 'fa-spinner fa-spin-pulse' : 'fa-file-pdf'}`" />
            PDF generieren
          </button>
        </div>

        <div class="card mb-3">
          <div class="card-body">
            <p
              v-if="!selectedSet"
              class="mb-0"
            >
              Zum Anzeigen der Felder oben den Dokumententyp auswählen.
            </p>
            <!-- Field-generator -->
            <template
              v-for="field in docFields"
              :key="field.name"
            >
              <!-- String, multi-line -->
              <div
                v-if="field.type === 'string' && field.format === 'multiline'"
                class="mb-3"
              >
                <label :for="`field-${field.name}`">{{ field.description }}</label>
                <textarea
                  :id="`field-${field.name}`"
                  v-model="model[field.name]"
                  :class="`form-control ${fieldValidClass(field.name)}`"
                />
              </div>

              <!-- String, enum -->
              <div
                v-else-if="field.type === 'string' && field.enum"
                class="mb-3"
              >
                <label :for="`field-${field.name}`">{{ field.description }}</label>
                <select
                  :id="`field-${field.name}`"
                  v-model="model[field.name]"
                  :class="`form-select ${fieldValidClass(field.name)}`"
                >
                  <option
                    v-for="opt in field.enum"
                    :key="opt"
                    :value="opt"
                  >
                    {{ opt }}
                  </option>
                </select>
              </div>

              <!-- String, single-line -->
              <div
                v-else-if="field.type === 'string'"
                class="mb-3"
              >
                <label :for="`field-${field.name}`">{{ field.description }}</label>
                <input
                  :id="`field-${field.name}`"
                  v-model="model[field.name]"
                  type="text"
                  :class="`form-control ${fieldValidClass(field.name)}`"
                >
              </div>

              <!-- Boolean -->
              <div
                v-else-if="field.type === 'boolean'"
                class="form-check form-switch mb-3"
              >
                <input
                  :id="`field-${field.name}`"
                  v-model="model[field.name]"
                  type="checkbox"
                  class="form-check-input"
                >
                <label :for="`field-${field.name}`">{{ field.description }}</label>
              </div>

              <!-- Number -->
              <div
                v-else-if="field.type === 'number' || field.format === 'integer'"
                class="mb-3"
              >
                <label :for="`field-${field.name}`">{{ field.description }}</label>
                <input
                  :id="`field-${field.name}`"
                  v-model.number="model[field.name]"
                  type="number"
                  :class="`form-control ${fieldValidClass(field.name)}`"
                >
              </div>
            </template>
          </div>
        </div>

        <div
          id="extraFields"
          class="accordion"
        >
          <!-- Add addresses -->
          <div class="accordion-item">
            <h2 class="accordion-header">
              <button
                class="accordion-button collapsed"
                type="button"
                data-bs-toggle="collapse"
                data-bs-target="#accAddress"
              >
                Serienbrief-Adressen hinzufügen
              </button>
            </h2>
            <div
              id="accAddress"
              class="accordion-collapse collapse"
              data-bs-parent="#extraFields"
            >
              <div class="accordion-body">
                <div class="input-group">
                  <input
                    id="recipientfile"
                    ref="csvInput"
                    type="file"
                    class="form-control"
                    accept=".csv"
                    @change="readRecipients"
                  >
                  <button
                    class="btn btn-danger"
                    @click="clearRecipients"
                  >
                    <i class="fas fa-trash fa-fw" />
                  </button>
                </div>
                <p class="mb-0 form-text">
                  Erwartete Felder: <code>NACHNAME;VORNAME;STRASSE;HAUSNR;PLZ;ORT</code>, Trenner <code>;</code>, eine Zeile pro Adresse
                </p>
              </div>
            </div>
          </div>

          <!-- Template -->
          <div class="accordion-item">
            <h2 class="accordion-header">
              <button
                class="accordion-button collapsed"
                type="button"
                data-bs-toggle="collapse"
                data-bs-target="#loadTemplate"
              >
                Vorlage laden / speichern...
              </button>
            </h2>
            <div
              id="loadTemplate"
              class="accordion-collapse collapse"
              data-bs-parent="#extraFields"
            >
              <div class="accordion-body">
                <div class="input-group mb-3">
                  <input
                    id="templatefile"
                    ref="templateFileInput"
                    type="file"
                    class="form-control"
                    accept=".json"
                    @change="readTemplateFromFile"
                  >
                </div>
                <div class="">
                  <a
                    class="btn form-control btn-primary"
                    :href="templateContentURL"
                    :download="selectedSet"
                  >
                    <i class="fas fa-download fa-fw me-1" />
                    Speichern...
                  </a>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>

      <div class="col h-100">
        <iframe
          v-if="displayURL"
          :src="displayURL"
          class="w-100 h-100"
        />
        <div
          v-else
          class="h-100 w-100 d-flex justify-content-center align-items-center text-center"
        >
          <p>
            <i class="fas fa-file-pdf fs-1 mb-3" /><br>
            PDF generieren um eine Vorschau zu erhalten.
          </p>
        </div>
      </div>
    </div>
  </div>
</template>

<script lang="ts">
import { Base64 } from 'js-base64'
import { Collapse } from 'bootstrap'
import { defineComponent } from 'vue'

export default defineComponent({
  computed: {
    docFields(): any[] {
      if (!this.selectedSet || !this.sourceSets[this.selectedSet]) {
        return []
      }

      const set = this.sourceSets[this.selectedSet]

      return Object.entries(set.properties)
        .map((e: any[]) => ({
          ...e[1],
          name: e[0],
          required: set.required.includes(e[0]),
        }))
    },

    docTypes(): any[] {
      return Object.entries(this.sourceSets)
        .map((e: any[]) => ({ description: e[1].description, name: e[0] }))
        .sort((a: any, b: any) => a.description.localeCompare(b.description))
    },

    modelFieldValid(): Record<string, boolean> {
      const validity = {}

      if (!this.selectedSet || !this.docFields) {
        return validity
      }

      for (const field of this.docFields) {
        // By default assume the field to be valid
        validity[field.name] = true

        if (field.required && this.model[field.name] === this.defaultForType(field.type)) {
          validity[field.name] = false
          continue
        }

        if (field.pattern && !this.model[field.name].match(new RegExp(field.pattern))) {
          validity[field.name] = false
          continue
        }
      }

      return validity
    },

    modelValid(): boolean {
      return Object.entries(this.modelFieldValid)
        .filter(e => !e[1])
        .length === 0
    },

    templateContentURL(): string {
      const data = {
        fields: this.model,
        type: this.selectedSet,
      }
      return `data:application/json;base64,${Base64.encode(JSON.stringify(data))}`
    },
  },

  data() {
    return {
      displayURL: '',
      documentLoading: false,
      model: {} as any,
      modelPrefill: {} as any,
      recipients: null as null | string,
      selectedSet: '',
      sourceSets: {} as any,
    }
  },

  methods: {
    clearRecipients(): void {
      this.recipients = null
      this.$refs.csvInput.value = ''
    },

    defaultForType(t: string, defaultValue?: any): any {
      if (defaultValue) {
        return defaultValue
      }

      switch (t) {
      case 'boolean':
        return false
      case 'string':
        return ''
      case 'number':
        return 0
      case 'integer':
        return 0
      }
    },

    fetchSourceSets(): Promise<void> {
      return fetch('/api/sets', {
        credentials: 'include',
      })
        .then((resp: Response) => resp.json())
        .then((data: any) => {
          this.sourceSets = data
        })
    },

    fieldValidClass(fieldName: string): string {
      return this.modelFieldValid[fieldName] ? '' : 'is-invalid'
    },

    loadTemplate(src: any): void {
      this.selectedSet = src.type
      this.modelPrefill = src.fields
    },

    readRecipients(): void {
      if ((this.$refs.csvInput.files?.length || 0) < 1) {
        this.recipients = null
        return
      }

      const file = this.$refs.csvInput.files[0] as File
      file.text()
        .then((csvContent: string) => {
          this.recipients = csvContent
        })
    },

    readTemplateFromFile(): void {
      if ((this.$refs.templateFileInput?.files?.length || 0) < 1) {
        console.warn(`loading template without files`)
        return
      }

      const file = this.$refs.templateFileInput.files[0] as File
      file.text()
        .then((content: string) => JSON.parse(content))
        .then((template: any) => this.loadTemplate(template))
    },

    readTemplateFromURL(url: string): void {
      fetch(url)
        .then((resp: Response) => resp.json())
        .then((template: any) => this.loadTemplate(template))
    },

    render(): Promise<void> | undefined {
      if (!this.modelValid) {
        return
      }

      this.documentLoading = true
      return fetch(`/api/render/${this.selectedSet}`, {
        body: JSON.stringify({
          foxCSV: this.recipients ? this.recipients : undefined,
          values: this.model,
        }),
        credentials: 'include',
        headers: {
          'Content-Type': 'application/json',
        },
        method: 'POST',
      })
        .then((resp: Response) => resp.blob())
        .then((data: Blob) => {
          this.displayURL = URL.createObjectURL(data)
          this.documentLoading = false
        })
    },
  },

  mounted(): void {
    this.fetchSourceSets()

    const collapseElementList = document.querySelectorAll('.collapse')
    for (const collapseEl of collapseElementList) {
      new Collapse(collapseEl, { toggle: false })
    }

    const hashParams = new URLSearchParams(window.location.hash.substring(1))
    if (hashParams.has('tplsrc')) {
      // Load after giving a tiny bit of time for the watcher not to escalate
      window.setTimeout(() => this.readTemplateFromURL(hashParams.get('tplsrc')), 100)
    }
  },

  name: 'DocRenderApp',

  watch: {
    selectedSet(to) {
      const fields = this.sourceSets[to].properties
      const model = {}

      for (const field of Object.entries(fields) as Array<Array<any>>) {
        model[field[0]] = this.modelPrefill[field[0]] ? this.modelPrefill[field[0]] : this.defaultForType(field[1].type, field[1].default)
      }

      this.model = model
    },
  },
})
</script>

<style scoped>
textarea.form-control {
  min-height: 200px;
}
</style>
