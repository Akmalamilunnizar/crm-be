export const useApiHost = () => {
  const config = useRuntimeConfig();
  let api = config.public.API_HOST;

  if (!api || api.trim() === '') {
    console.warn('⚠️ NUXT_PUBLIC_API_HOST is not defined in .env; using fallback http://localhost:3001');
    api = 'http://localhost:3001';
  }

  return api;
};
export const useWaHost = () => {
  const config = useRuntimeConfig();
  let wa = config.public.WA_HOST;

  if (!wa || wa.trim() === '') { 
    console.warn('⚠️ NUXT_PUBLIC_API_HOST is not defined in .env; using fallback http://localhost:3001');
    wa = 'http://localhost:3001';
  }

  return wa;
};