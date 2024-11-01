import { expect, test, describe, beforeAll, afterAll, afterEach, vi } from 'vitest';
import { http, HttpResponse } from 'msw';
import { load } from './+page';
import { render, screen } from '@testing-library/svelte';
import { setupServer } from 'msw/node';
import { BACKEND_SERVER } from '$lib/constants';
import CarPage from './+page.svelte';
import { test_data } from './test_data';
import userEvent from '@testing-library/user-event';
import type { PageData, PageLoadEvent } from './$types';
import { mock } from 'vitest-mock-extended';

const server = setupServer();
const user = userEvent.setup();

beforeAll(() => {
    // NOTE: server.listen must be called before `createClient` is used to ensure
    // the msw can inject its version of `fetch` to intercept the requests.
    server.listen({
        onUnhandledRequest: (request) => {
            throw new Error(`No request handler found for ${request.method} ${request.url}`);
        }
    });
});

afterEach(() => server.resetHandlers());

afterAll(() => server.close());

describe('fetch cars information test', () => {
    test('test if cars is loaded correctly', async () => {
        //Data is loaded correctly
        server.use(
            http.get(`${BACKEND_SERVER}/cars`, () => HttpResponse.json(test_data, { status: 200 }))
        );

        const loadEvent = mock<PageLoadEvent>({ fetch: global.fetch });
        const data = (await load(loadEvent)) as PageData;
        expect(data.cars).toStrictEqual(test_data);

        render(CarPage, { data: data });
        test_data.forEach((car) => {
            screen.getByText(car.details.license_plate);
            screen.getByText('Color: ' + car.details.color);
            screen.getByText('Make: ' + car.details.make);
            screen.getByText('Model: ' + car.details.model);
        });
    });

    test('test if cars create work correctly', async () => {
        const data = mock<PageData>({ cars: test_data });
        const new_car_detail = {
            license_plate: 'lic-new',
            make: 'color-new',
            model: 'model-new',
            color: 'make-new'
        };
        render(CarPage, { data: data });

        server.use(
            http.post(`${BACKEND_SERVER}/cars`, () =>
                HttpResponse.json({ details: new_car_detail, id: 'random-id' }, { status: 201 })
            )
        );

        const createButton = screen.getByRole('button', { name: 'new-car-button' });
        await user.click(createButton);
        let license_plate_field = screen.getByRole('textbox', {
            name: 'License plate'
        });
        let color_field = screen.getByRole('textbox', { name: 'Color' });
        let model_field = screen.getByRole('textbox', { name: 'Model' });
        let make_field = screen.getByRole('textbox', { name: 'Make' });

        console.log('Test creating a new car with correct input');
        await user.click(license_plate_field);
        await user.keyboard(new_car_detail.license_plate);

        await user.click(color_field);
        await user.keyboard(new_car_detail.color);

        await user.click(model_field);
        await user.keyboard(new_car_detail.model);

        await user.click(make_field);
        await user.keyboard(new_car_detail.make);

        let confirmButton = screen.getByRole('button', { name: 'Confirm' });
        await user.click(confirmButton);
        screen.getByText(new_car_detail.license_plate);
        screen.getByText('Color: ' + new_car_detail.color);
        screen.getByText('Make: ' + new_car_detail.make);
        screen.getByText('Model: ' + new_car_detail.model);

        server.use(
            http.post(`${BACKEND_SERVER}/cars`, () =>
                HttpResponse.json(
                    {
                        details: 'the specified license plate is invalid',
                        errors: [
                            {
                                location: 'body.license_plate',
                                value: 'very wrong license'
                            }
                        ]
                    },
                    { status: 422 }
                )
            )
        );
        console.log('Test creating a new car with incorrect input');
        await user.click(createButton);

        license_plate_field = screen.getByRole('textbox', {
            name: 'License plate'
        });
        color_field = screen.getByRole('textbox', { name: 'Color' });
        model_field = screen.getByRole('textbox', { name: 'Model' });
        make_field = screen.getByRole('textbox', { name: 'Make' });
        confirmButton = screen.getByRole('button', { name: 'Confirm' });
        await user.click(license_plate_field);
        await user.keyboard('Very wrong license this is certainly invalid');

        await user.click(color_field);
        await user.keyboard('wrong');

        await user.click(model_field);
        await user.keyboard('wrong');

        await user.click(make_field);
        await user.keyboard('wrong');

        await user.click(confirmButton);
        screen.getByText('Error:');
        const closeModalButton = screen.getByRole('button', { name: 'Close the modal' });
        await user.click(closeModalButton);
    });

    test('Test delete functionality', async () => {
        const data = mock<PageData>({ cars: test_data });
        window.confirm = vi.fn(() => {
            console.log('confirm');
            return true;
        });

        render(CarPage, { data: data });

        server.use(
            http.delete(`${BACKEND_SERVER}/cars/:id`, () =>
                HttpResponse.json({ data: 'random' }, { status: 204 })
            )
        );
        const car_to_delete = screen.getByText(test_data[0].details.license_plate);
        await user.click(car_to_delete);
        const deleteButton = screen.getByRole('button', { name: 'Delete' });
        await user.click(deleteButton);
        const deletedCar = screen.queryByText(test_data[0].details.license_plate);
        expect(deletedCar).toBe(null);
    });

    test('Test edit functionality', async () => {
        const data = mock<PageData>({ cars: test_data });
        render(CarPage, { data: data });
        const car_to_edit = screen.getByText(test_data[0].details.license_plate);
        await user.click(car_to_edit);
        let editButton = screen.getByRole('button', { name: 'Edit' });
        await user.click(editButton);

        const edit_car_detail = {
            license_plate: 'lic-edit',
            make: 'color-edit',
            model: 'model-edit',
            color: 'make-edit'
        };

        server.use(
            http.put(`${BACKEND_SERVER}/cars/:id`, () =>
                HttpResponse.json(
                    { details: edit_car_detail, id: test_data[0].id },
                    { status: 200 }
                )
            )
        );
        console.log('Edit with correct input');
        let license_plate_field = screen.getByRole('textbox', { name: 'License plate' });
        let color_field = screen.getByRole('textbox', { name: 'Color' });
        let model_field = screen.getByRole('textbox', { name: 'Model' });
        let make_field = screen.getByRole('textbox', { name: 'Make' });
        await user.click(license_plate_field);
        await user.clear(license_plate_field);
        await user.keyboard(edit_car_detail.license_plate);

        await user.click(color_field);
        await user.clear(color_field);
        await user.keyboard(edit_car_detail.color);

        await user.click(model_field);
        await user.clear(model_field);
        await user.keyboard(edit_car_detail.model);

        await user.click(make_field);
        await user.clear(make_field);
        await user.keyboard(edit_car_detail.make);

        let confirmButton = screen.getByRole('button', { name: 'Confirm' });
        await user.click(confirmButton);
        screen.getAllByText(edit_car_detail.license_plate);
        screen.getAllByText('Color: ' + edit_car_detail.color);
        screen.getAllByText('Make: ' + edit_car_detail.make);
        screen.getAllByText('Model: ' + edit_car_detail.model);

        console.log('Edit with wrong input');
        server.use(
            http.put(`${BACKEND_SERVER}/cars/:id`, () =>
                HttpResponse.json(
                    {
                        details: 'the specified license plate is invalid',
                        errors: [
                            {
                                location: 'body.license_plate',
                                value: 'very wrong license'
                            }
                        ]
                    },
                    { status: 422 }
                )
            )
        );

        editButton = screen.getByRole('button', { name: 'Edit' });
        await user.click(editButton);
        license_plate_field = screen.getByRole('textbox', { name: 'License plate' });
        color_field = screen.getByRole('textbox', { name: 'Color' });
        model_field = screen.getByRole('textbox', { name: 'Model' });
        make_field = screen.getByRole('textbox', { name: 'Make' });
        await user.click(license_plate_field);
        await user.clear(license_plate_field);
        await user.keyboard('1');

        await user.click(color_field);
        await user.clear(color_field);
        await user.keyboard('1');

        await user.click(model_field);
        await user.clear(model_field);
        await user.keyboard('1');

        await user.click(make_field);
        await user.clear(make_field);
        await user.keyboard('1');
        confirmButton = screen.getByRole('button', { name: 'Confirm' });

        await user.click(confirmButton);

        screen.getByText('Error:');
        const closeModalButton = screen.getByRole('button', { name: 'Close' });
        await user.click(closeModalButton);
    });
});
