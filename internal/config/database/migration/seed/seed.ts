import { v4 as uuidv4 } from 'uuid';
import { PrismaClient } from '../generated/client';
const prisma = new PrismaClient()
async function main() {
    const roles = prisma.roles.createMany({
        data: [
            {
                name: 'Admin',
                id: uuidv4()
            },
            {
                name: 'User',
                id: uuidv4()
            }
        ],
        skipDuplicates: true
    })
}
main()
    .then(async () => {
        await prisma.$disconnect()
    })
    .catch(async (e) => {
        console.error(e)
        await prisma.$disconnect()
    })