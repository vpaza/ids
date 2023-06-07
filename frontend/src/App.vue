
<template>
  <div class="flex flex-col min-h-screen">
    <header class="sticky z-50 bg-neutral-700 text-white top-0 p-4">
      ZAN IDS
    </header>
    <div class="flex-grow">
      <main>
        <router-view />
      </main>
    </div>
    <footer class="sticky z-50 bg-neutral-700 text-white bottom-0 p-4">
      footer
    </footer>
  </div>
</template>

<style scoped>

</style>
